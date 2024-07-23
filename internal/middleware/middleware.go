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

	"github.com/gleb-korostelev/short-url/internal/config"
	"github.com/gleb-korostelev/short-url/internal/service/utils"
	"github.com/gleb-korostelev/short-url/tools/logger"
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

// WriteHeader sends an HTTP response header with the status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Multiple calls to WriteHeader will result in an invalid multiple header write error.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode                    // Store the status code to handle multiple calls or checks.
	w.ResponseWriter.WriteHeader(statusCode) // Delegate to the underlying ResponseWriter.
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not been called, it calls WriteHeader(http.StatusOK)
// before writing the data. It returns the number of bytes written and any write error encountered.
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 { // Check if the status code has not been set.
		w.WriteHeader(http.StatusOK) // Set the default status code to OK if not set.
	}
	size, err := w.ResponseWriter.Write(b) // Write the data using the embedded ResponseWriter.
	w.size += size                         // Update the size of the data written.
	return size, err                       // Return the size of the data written and any error encountered.
}

// gzipResponseWriter is an enhanced http.ResponseWriter that supports Gzip compression.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	hasWrittenHeader bool
}

// WriteHeader sets the status code for the HTTP response header.
// If WriteHeader is called after writing has started, it returns without modifying the header.
// This prevents headers from being rewritten which can lead to protocol errors.
func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.hasWrittenHeader {
		return // Prevent modification after headers are written.
	}
	w.hasWrittenHeader = true
	w.ResponseWriter.WriteHeader(statusCode) // Delegate to the underlying ResponseWriter to set the status code.
}

// Write writes the provided byte slice into the response body.
// If WriteHeader has not yet been called, Write will first set the Content-Encoding to gzip
// for specific content types ('application/json', 'text/html') and initialize gzip compression.
// If the content type is not compatible with gzip, it writes directly using the underlying ResponseWriter.
// It ensures headers are written before any response body if not already done.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.hasWrittenHeader {
		// Automatically handle content encoding and compression based on content type
		contentType := w.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w.ResponseWriter)
			defer gz.Close() // Ensure the gzip writer is closed after the write operation
			w.Writer = gz    // Use gzip writer for response body
		} else {
			w.Writer = w.ResponseWriter // Use the normal response writer for non-compatible types
		}
		w.WriteHeader(http.StatusOK) // Set default status code if not yet set
	}
	return w.Writer.Write(b) // Write the data to the selected writer
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
