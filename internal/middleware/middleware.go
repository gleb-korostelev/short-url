// Package middleware contains HTTP middleware functions that enhance the HTTP server functionality.
// These middlewares provide logging, compression, decompression, and user authentication handling.
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

// LoggingMiddleware logs the HTTP request details including method, URI, status code, response size, and duration.
// It uses the zap.Logger for structured logging.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(ww, r)
			go logger.Info("request",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", ww.status),
				zap.Int("response_size", ww.size),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

// GzipCompressMiddleware applies Gzip compression to HTTP responses if the client supports it.
// It checks the Accept-Encoding header for 'gzip' and wraps the response writer to compress the output.
func GzipCompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: w}, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// GzipDecompressMiddleware handles Gzip-compressed request bodies.
// It checks the Content-Encoding header for 'gzip' and decompresses the body if necessary.
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

// responseWriter is a custom http.ResponseWriter that captures HTTP status codes and response sizes for logging.
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

// gzipResponseWriter is an enhanced http.ResponseWriter that supports Gzip compression.
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
			gz := gzip.NewWriter(w.ResponseWriter)
			defer gz.Close()
			w.Writer = gz
		} else {
			w.Writer = w.ResponseWriter
		}
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}

// EnsureUserCookie checks for a valid user ID from cookies.
// If not found, it generates a new user ID, sets it in a cookie, and logs the error.
// It authorizes users by ensuring a valid user ID is present or created.
func EnsureUserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromCookie(r)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || err == config.ErrTokenInvalid {
				userID = uuid.New().String()
				utils.SetJWTInCookie(w, userID)
				logger.Infof("Generated new user ID and set in cookie due to error: %v", err)
			} else {
				logger.Infof("Failed to authorize due to error: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		ctx := context.WithValue(r.Context(), config.UserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
