// Package storage defines an interface for URL storage operations. It abstracts the underlying
// storage mechanism, which could be in-memory, file-based, or database-backed. Implementations
// of this interface handle saving, retrieving, and managing URLs.
package storage

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/models"
)

// Storage is the interface that defines the methods required to store, retrieve, and manage URLs.
// Implementations of this interface must handle various storage operations, including URL creation,
// retrieval, and lifecycle management in a thread-safe manner.
type Storage interface {
	// SaveUniqueURL stores a new URL and associates it with a user ID, ensuring the short URL is unique.
	// Returns the shortened URL, an HTTP status code indicating the result, and any error encountered.
	SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error)

	// SaveURL stores a new URL without ensuring uniqueness.
	// It is typically used when the unique handling is managed at a higher level or not required.
	SaveURL(ctx context.Context, originalURL string, userID string) (string, error)

	// GetOriginalLink retrieves the original URL based on its shortened version.
	// It returns the original URL and any error encountered if the URL does not exist or other issues arise.
	GetOriginalLink(ctx context.Context, shortURL string) (string, error)

	// Ping checks the health or connectivity of the storage medium, often used in database connections.
	// It returns an HTTP status code and any error encountered during the health check.
	Ping(ctx context.Context) (int, error)

	// Close performs cleanup or closure operations on the storage, such as closing database connections.
	Close() error

	// GetAllURLS retrieves all URLs associated with a specific user ID.
	// This method is useful for user-specific URL management and returns a slice of UserURLs and any error encountered.
	GetAllURLS(ctx context.Context, userID, baseURL string) ([]models.UserURLs, error)

	// MarkURLsAsDeleted marks specified URLs as deleted for a given user ID.
	// This method handles the soft deletion of URLs and returns any error encountered during the operation.
	MarkURLsAsDeleted(ctx context.Context, userID string, shortURLs []string) error
}
