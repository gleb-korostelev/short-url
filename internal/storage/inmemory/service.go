package inmemory

import (
	"context"
	"net/http"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

type service struct {
	cache map[string]string
	mu    *sync.RWMutex
}

func NewMemoryStorage(cache map[string]string, mut *sync.RWMutex) storage.Storage {
	return &service{
		cache: cache,
		mu:    mut,
	}
}

func (s *service) SaveUniqueURL(ctx context.Context, originalURL string) (string, int, error) {
	shortURL := business.GenerateShortPath()
	for _, exists := s.cache[shortURL]; exists; {
		shortURL = business.GenerateShortPath()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL, http.StatusCreated, nil
}

func (s *service) SaveURL(ctx context.Context, originalURL string) (string, error) {
	shortURL := business.GenerateShortPath()
	for _, exists := s.cache[shortURL]; exists; {
		shortURL = business.GenerateShortPath()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL, nil
}

func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	originalURL, exists := s.cache[shortURL]
	if !exists {
		return "", config.ErrNotFound
	}
	return originalURL, nil
}

func (s *service) Ping(ctx context.Context) (int, error) {
	logger.Errorf("Using inmemory save %v", config.ErrWrongMode)
	return http.StatusInternalServerError, config.ErrWrongMode
}

func (s *service) Close() error {
	return nil
}
