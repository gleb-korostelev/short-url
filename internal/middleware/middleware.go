package middleware

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(ww, r)

			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", ww.status),
				zap.Int("response_size", ww.size),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func GzipCompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: w}, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func GzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusBadRequest)
				return
			}
			defer gzReader.Close()
			r.Body = gzReader
			r.Header.Del("Content-Encoding")
		}
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	hasWrittenHeader bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.hasWrittenHeader {
		return
	}
	w.hasWrittenHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.hasWrittenHeader {
		contentType := w.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Writer = gzip.NewWriter(w.ResponseWriter)
			defer func() {
				if w.Writer != nil {
					w.Writer.(*gzip.Writer).Close()
				}
			}()
		}
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}

func EnsureUserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromCookie(r)
		if errors.Is(err, http.ErrNoCookie) || err == config.ErrTokenInvalid {
			userID = uuid.New().String()
			utils.SetJWTInCookie(w, userID)
			logger.Infof("error in cookie is %v", err)
		} else if err != nil {
			logger.Infof("error in cookie auth is %v", err)
			http.Error(w, "Failed to Authorize", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), config.UserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
