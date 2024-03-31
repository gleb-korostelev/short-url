package inmemory

import (
	"context"
	"net/http"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

type service struct {
	cache []models.URLData
	mu    sync.RWMutex
}

func NewMemoryStorage(cache []models.URLData) storage.Storage {
	return &service{
		cache: cache,
	}
}

func (s *service) SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userId in database %v", err)
		return "", http.StatusBadRequest, err
	}

	shortURL := business.GenerateShortPath()

	for _, info := range s.cache {
		if info.ShortURL == shortURL {
			shortURL = business.GenerateShortPath()
		}
	}

	var data models.URLData
	data.ShortURL = shortURL
	data.OriginalURL = originalURL
	data.UUID = uuid
	data.DeletedFlag = false
	s.cache = append(s.cache, data)

	s.mu.RLock()
	defer s.mu.RUnlock()
	return config.BaseURL + "/" + shortURL, http.StatusCreated, nil
}

func (s *service) SaveURL(ctx context.Context, originalURL string, userID string) (string, error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userId in database %v", err)
		return "", err
	}

	shortURL := business.GenerateShortPath()

	for _, info := range s.cache {
		if info.ShortURL == shortURL {
			shortURL = business.GenerateShortPath()
		}
	}

	var data models.URLData
	data.ShortURL = shortURL
	data.OriginalURL = originalURL
	data.UUID = uuid
	data.DeletedFlag = false
	s.cache = append(s.cache, data)

	s.mu.RLock()
	defer s.mu.RUnlock()
	return config.BaseURL + "/" + shortURL, nil
}

func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, info := range s.cache {
		if info.ShortURL == shortURL {
			if info.DeletedFlag {
				return "", config.ErrGone
			}
			return info.OriginalURL, nil
		}
	}
	return "", config.ErrNotFound

}

func (s *service) Ping(ctx context.Context) (int, error) {
	logger.Errorf("Using inmemory save %v", config.ErrWrongMode)
	return http.StatusInternalServerError, config.ErrWrongMode
}

func (s *service) Close() error {
	return nil
}

func (s *service) GetAllURLS(ctx context.Context, userID string) ([]models.AllUserURL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var urls []models.AllUserURL
	var data models.AllUserURL
	for _, info := range s.cache {
		if info.UUID.String() == userID && !info.DeletedFlag {
			data.OriginalURL = info.OriginalURL
			data.ShortURL = config.BaseURL + "/" + info.ShortURL
			urls = append(urls, data)
		}
	}
	return urls, nil
}

func (s *service) MarkURLsAsDeleted(ctx context.Context, userID string, shortURLs []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, info := range s.cache {
		if info.UUID.String() == userID && business.CheckURL(info.ShortURL, shortURLs) {
			s.cache[i].DeletedFlag = true
		}

	}

	return nil
}
