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

type service struct {
	cache map[string]models.URLData
	mu    sync.RWMutex
}

func NewMemoryStorage(cache map[string]models.URLData) storage.Storage {
	return &service{
		cache: cache,
	}
}

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

func (s *service) Ping(ctx context.Context) (int, error) {
	logger.Errorf("Using inmemory save %v", config.ErrWrongMode)
	return http.StatusInternalServerError, config.ErrWrongMode
}

func (s *service) Close() error {
	return nil
}

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
