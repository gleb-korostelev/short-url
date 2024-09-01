package grpcservice

import (
	"context"
	"errors"
	"net"
	"net/http"

	pb "github.com/gleb-korostelev/short-url.git/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
)

// URLServiceServer should implement all the server APIs for URL service,
type URLServiceServerImpl struct {
	pb.UnimplementedURLServiceServer
	store  storage.Storage
	worker worker.DBWorkerPool
}

// GetOriginal retrieves the original URL corresponding to a given shortened URL ID.
func (s *URLServiceServerImpl) GetOriginal(ctx context.Context, req *pb.GetOriginalRequest) (*pb.GetOriginalResponse, error) {
	shortURL := req.GetShortUrlId()

	originalURL, err := s.store.GetOriginalLink(ctx, shortURL)
	if err != nil {
		if errors.Is(err, config.ErrGone) {
			return nil, status.Errorf(codes.NotFound, config.ErrGone.Error(), shortURL)
		}
		return nil, status.Errorf(codes.Internal, "Error retrieving URL: %v", err)
	}

	return &pb.GetOriginalResponse{OriginalUrl: originalURL}, nil
}

// Ping handles healthcheck of database
func (s *URLServiceServerImpl) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	httpStatusCode, err := s.store.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(httpToGRPCStatusCode(httpStatusCode), "Failed to ping database: %v", err)
	}
	return &pb.PingResponse{Status: int32(httpStatusCode)}, nil
}

// PostShorter processes a request to create a shortened URL from an original URL.
func (s *URLServiceServerImpl) PostShorter(ctx context.Context, req *pb.PostShorterRequest) (*pb.PostShorterResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	originalURL := req.Url
	resultChan := make(chan *pb.PostShorterResponse)
	errorChan := make(chan error)

	s.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			shortURL, httpStatusCode, err := s.store.SaveUniqueURL(ctx, originalURL, userID)
			if err != nil {
				errorChan <- status.Errorf(httpToGRPCStatusCode(httpStatusCode), "failed to save URL: %v", err)
				return err
			}
			resultChan <- &pb.PostShorterResponse{ShortUrl: shortURL}
			return nil
		},
	})

	select {
	case res := <-resultChan:
		return res, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, status.Errorf(codes.DeadlineExceeded, "request timed out")
	}
}

// PostShorterJSON handles gRPC requests to create shortened URLs from a provided original URL in JSON format.
func (s *URLServiceServerImpl) PostShorterJSON(ctx context.Context, req *pb.PostShorterJSONRequest) (*pb.PostShorterJSONResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	originalURL := req.Url
	if originalURL == "" {
		return nil, status.Error(codes.InvalidArgument, "original URL cannot be empty")
	}

	shortURL, httpStatusCode, err := s.store.SaveUniqueURL(ctx, originalURL, userID)
	if err != nil {
		return nil, status.Errorf(httpToGRPCStatusCode(httpStatusCode), "failed to save URL: %v", err)
	}

	return &pb.PostShorterJSONResponse{
		ShortUrl: shortURL,
	}, nil
}

// ShortenBatch processes multiple URL shorten requests simultaneously.
func (s *URLServiceServerImpl) ShortenBatch(ctx context.Context, req *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	if len(req.UrlPairs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Empty batch is not allowed")
	}

	var results []*pb.ShortURLPairResp
	for _, item := range req.UrlPairs {
		shortURL, err := s.store.SaveURL(ctx, item.OriginalUrl, userID)
		if err != nil {
			return &pb.ShortenBatchResponse{UrlPairs: results}, status.Errorf(codes.Internal, "Error with saving URL: %v", err)
		}
		results = append(results, &pb.ShortURLPairResp{
			CorrelationId: item.CorrelationId,
			ShortUrl:      shortURL,
		})
	}

	return &pb.ShortenBatchResponse{UrlPairs: results}, nil
}

// GetUserURLs retrieves all URLs associated with the authenticated user.
func (s *URLServiceServerImpl) GetUserURLs(ctx context.Context, req *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	urls, err := s.store.GetAllURLS(ctx, userID, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve URLs: %v", err)
	}

	if len(urls) == 0 {
		return &pb.GetUserURLsResponse{}, status.Error(codes.NotFound, "no URLs found for the user")
	}

	userURLs := make([]*pb.UserURL, len(urls))
	for i, url := range urls {
		userURLs[i] = &pb.UserURL{
			ShortUrl:    url.ShortURL,
			OriginalUrl: url.OriginalURL,
		}
	}

	return &pb.GetUserURLsResponse{UserUrls: userURLs}, nil
}

// DeleteURLs marks URLs as deleted based on the provided request and current user authentication.
func (s *URLServiceServerImpl) DeleteURLs(ctx context.Context, req *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	shortURLs := req.GetShortUrls()
	if len(shortURLs) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no URLs provided for deletion")
	}

	doneChan := make(chan error)
	s.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			err := s.store.MarkURLsAsDeleted(ctx, userID, shortURLs)
			doneChan <- err
			return err
		},
	})

	err = <-doneChan
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting URLs: %v", err)
	}

	return &pb.DeleteURLsResponse{Status: int32(codes.OK)}, nil
}

// Stats provides statistics about the service, such as the number of shortened URLs and registered users.
func (s *URLServiceServerImpl) Stats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	if config.TrustedSubnet == "" {
		return nil, status.Errorf(codes.PermissionDenied, "access denied due to missing trusted subnet configuration")
	}

	clientIP, err := extractClientIP(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract client IP: %v", err)
	}

	if !isIPAllowed(clientIP, config.TrustedSubnet) {
		return nil, status.Errorf(codes.PermissionDenied, "access from IP %s is not allowed", clientIP)
	}

	urlsCount, usersCount, err := s.store.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error fetching stats: %v", err)
	}

	return &pb.StatsResponse{
		UrlsCount:  int32(urlsCount),
		UsersCount: int32(usersCount),
	}, nil
}

// extractClientIP extracts the client's IP address from gRPC metadata.
func extractClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", status.Error(codes.Internal, "failed to get client peer info")
	}

	return p.Addr.String(), nil
}

// isIPAllowed checks if the given IP address is within the trusted subnet.
func isIPAllowed(ip, trustedSubnet string) bool {
	_, subnet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false
	}

	return subnet.Contains(net.ParseIP(ip))
}

// getUserIDFromContext gets UserKey from context
func getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(config.UserContextKey).(string)
	if !ok {
		return "", config.ErrUserNotFound
	}
	return userID, nil
}

// httpToGRPCStatusCode converts HTTP status codes to gRPC status codes.
func httpToGRPCStatusCode(httpStatusCode int) codes.Code {
	switch httpStatusCode {
	case http.StatusOK:
		return codes.OK
	case http.StatusInternalServerError:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
