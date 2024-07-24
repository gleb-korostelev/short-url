// Package inmemory implements the storage.Storage interface using an in-memory data store.
// It provides fast access to URL data stored directly in memory, suitable for scenarios
// where persistence across service restarts is not required.
package inmemory

import (
	"context"
	"net/http"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

// service provides an in-memory storage mechanism for URL data.
// It uses a map to store URL data, keyed by short URL strings, and a mutex to manage concurrent access.
type service struct {
	cache map[string]models.URLData // cache stores the URL data in-memory.
	mu    sync.RWMutex              // mu protects the cache from concurrent read/write access.
}

// NewMemoryStorage initializes a new in-memory storage service with a given initial cache.
func NewMemoryStorage(cache map[string]models.URLData) storage.Storage {
	return &service{
		cache: cache,
	}
}

// SaveUniqueURL saves a new URL into the in-memory storage, ensuring the short URL is unique.
// It generates a short URL, checks for uniqueness within the existing entries, and saves the URL data.
// Returns the complete URL, HTTP status code, and error if any.
func (s *service) SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userId in memory %v", err)
		return "", http.StatusBadRequest, err
	}

	shortURL := utils.GenerateShortPath()

	for _, exists := s.cache[shortURL]; exists; {
		shortURL = utils.GenerateShortPath()
	}

	var data models.URLData
	data.ShortURL = shortURL
	data.OriginalURL = originalURL
	data.UUID = uuid
	data.DeletedFlag = false
	s.cache[data.ShortURL] = data

	s.mu.RLock()
	defer s.mu.RUnlock()
	return config.BaseURL + "/" + shortURL, http.StatusCreated, nil
}

// SaveURL performs a similar operation to SaveUniqueURL but does not return an HTTP status.
func (s *service) SaveURL(ctx context.Context, originalURL string, userID string) (string, error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userId in memory %v", err)
		return "", err
	}

	shortURL := utils.GenerateShortPath()

	for _, exists := s.cache[shortURL]; exists; {
		shortURL = utils.GenerateShortPath()
	}

	var data models.URLData
	data.ShortURL = shortURL
	data.OriginalURL = originalURL
	data.UUID = uuid
	data.DeletedFlag = false
	s.cache[data.ShortURL] = data

	s.mu.RLock()
	defer s.mu.RUnlock()
	return config.BaseURL + "/" + shortURL, nil
}

// GetOriginalLink retrieves the original URL from a given short URL, checking if it's marked as deleted.
func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	foundCache, exists := s.cache[shortURL]
	if !exists {
		return "", config.ErrNotFound
	}
	if foundCache.DeletedFlag {
		return "", config.ErrGone
	}
	return foundCache.OriginalURL, nil

}

// Ping checks the operation status of the in-memory storage, typically returning an error as it does not involve connectivity.
func (s *service) Ping(ctx context.Context) (int, error) {
	logger.Errorf("Using inmemory save %v", config.ErrWrongMode)
	return http.StatusInternalServerError, config.ErrWrongMode
}

// Close performs cleanup if necessary; in this implementation, it is a no-operation.
func (s *service) Close() error {
	return nil
}

// GetAllURLs retrieves all URLs associated with a specific user ID, filtering out deleted entries.
func (s *service) GetAllURLS(ctx context.Context, userID, baseURL string) ([]models.UserURLs, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var urls []models.UserURLs
	var data models.UserURLs
	for _, info := range s.cache {
		if info.UUID.String() == userID && !info.DeletedFlag {
			data.OriginalURL = info.OriginalURL
			data.ShortURL = baseURL + "/" + info.ShortURL
			urls = append(urls, data)
		}
	}
	return urls, nil
}

// MarkURLsAsDeleted marks specified URLs as deleted for a given user ID.
func (s *service) MarkURLsAsDeleted(ctx context.Context, userID string, shortURLs []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, info := range s.cache {
		if info.UUID.String() == userID && utils.CheckURL(info.ShortURL, shortURLs) {
			info.DeletedFlag = true
			s.cache[info.ShortURL] = info
		}

	}

	return nil
}
