// Package filecache implements the storage.Storage interface using a file-based system
// to manage URL data. This package allows storing, retrieving, and deleting URLs directly on the filesystem.
package filecache

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

// service implements the storage.Storage interface to provide file-based URL management.
type service struct {
	path string // path represents the file path where URL data is stored.
}

// NewFileStorage creates a new instance of a file-based storage service.
// It accepts a file path where URL data will be stored and manipulated.
func NewFileStorage(path string) storage.Storage {
	return &service{
		path: path,
	}
}

// SaveUniqueURL saves a URL to the file and generates a unique short URL.
// It returns the created short URL, an HTTP status code, and any error encountered.
func (s *service) SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error) {
	shortURL := utils.GenerateShortPath()

	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userID in file %v", err)
		return "", http.StatusBadRequest, err
	}

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = uuid
	save.DeletedFlag = false

	err = utils.SaveURLs(save)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	return config.BaseURL + "/" + shortURL, http.StatusCreated, nil
}

// SaveURL saves a URL without ensuring uniqueness.
func (s *service) SaveURL(ctx context.Context, originalURL string, userID string) (string, error) {
	shortURL := utils.GenerateShortPath()

	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userID in file %v", err)
		return "", err
	}

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = uuid
	save.DeletedFlag = false

	err = utils.SaveURLs(save)
	if err != nil {
		logger.Errorf("Error with saving in file %v", err)
		return "", err
	}
	return config.BaseURL + "/" + shortURL, nil
}

// GetOriginalLink retrieves the original URL from the file for a given short URL.
func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := utils.LoadURLs(s.path, shortURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

// Ping simulates a connectivity test to the storage. Since this is a file-based system,
// the function returns an error indicating that this is a non-database mode.
func (s *service) Ping(ctx context.Context) (int, error) {
	logger.Errorf("Using file save %v", config.ErrWrongMode)
	return http.StatusInternalServerError, config.ErrWrongMode
}

// Close is a placeholder to satisfy the storage.Storage interface. It performs no operation.
func (s *service) Close() error {
	return nil
}

// GetAllURLS retrieves all URLs associated with a user ID.
func (s *service) GetAllURLS(ctx context.Context, userID, baseURL string) ([]models.UserURLs, error) {
	res, err := utils.LoadUserURLs(config.BaseFilePath, userID)
	if err != nil {
		logger.Errorf("Failed to get all user URLs %v", err)
		return nil, err
	}
	return res, nil
}

// MarkURLsAsDeleted marks specified URLs as deleted in the file system for a given user ID.
func (s *service) MarkURLsAsDeleted(ctx context.Context, userID string, shortURLs []string) error {
	err := utils.MarkURLsAsDeletedInFile(config.BaseFilePath, userID, shortURLs)
	return err
}
