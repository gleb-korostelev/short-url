package filecache

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

type service struct {
	path string
}

func NewFileStorage(path string) storage.Storage {
	return &service{
		path: path,
	}
}

func (s *service) SaveUniqueURL(ctx context.Context, originalURL string) (string, int, error) {

	shortURL := business.GenerateShortPath()

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = uuid.New()
	err := business.SaveURLs(save)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	return config.BaseURL + "/" + shortURL, http.StatusCreated, nil
}

func (s *service) SaveURL(ctx context.Context, originalURL string) (string, error) {
	shortURL := business.GenerateShortPath()

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = uuid.New()
	err := business.SaveURLs(save)
	if err != nil {
		logger.Errorf("Error with saving in file %v", err)
		return "", err
	}
	return config.BaseURL + "/" + shortURL, nil
}

func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := business.LoadURLs(s.path, shortURL)
	if err != nil {
		return "", err
	}
	return originalURL, err
}

func (s *service) Ping(ctx context.Context) (int, error) {
	logger.Errorf("Using file save %v", config.ErrWrongMode)
	return http.StatusInternalServerError, config.ErrWrongMode
}

func (s *service) Close() error {
	return nil
}
