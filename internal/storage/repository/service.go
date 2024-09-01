// Package repository implements the storage.Storage interface using a relational database.
// It provides comprehensive URL management functionalities including saving, retrieving,
// and deleting URLs through database transactions.
package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/db/dbimpl"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

// service provides URL storage management using a database.
type service struct {
	data db.DB // data is the interface for interacting with the database.
}

// NewDBStorage creates a new instance of a database-backed storage service.
func NewDBStorage(data db.DB) storage.Storage {
	return &service{
		data: data,
	}
}

// SaveUniqueURL saves a new URL into the database, ensuring it is unique.
// It generates a short URL and attempts to store it along with the original URL in the database.
// Returns the complete URL, HTTP status code, and any error encountered.
func (s *service) SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error) {
	shortURL := utils.GenerateShortPath()

	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error parsing userID: %v", err)
		return "", http.StatusInternalServerError, err
	}

	err = dbimpl.CreateShortURL(s.data, uuid.String(), shortURL, originalURL)
	if err != nil {
		if errors.Is(err, config.ErrExists) {
			existingShortURL, err := dbimpl.GetShortURLByOriginalURL(s.data, originalURL)
			if err != nil {
				return "", http.StatusInternalServerError, err
			}
			return config.BaseURL + "/" + existingShortURL, http.StatusConflict, nil
		}
		return "", http.StatusInternalServerError, err
	}
	return config.BaseURL + "/" + shortURL, http.StatusCreated, nil
}

// SaveURL performs the same operation as SaveUniqueURL without returning the HTTP status code.
func (s *service) SaveURL(ctx context.Context, originalURL string, userID string) (string, error) {
	shortURL := utils.GenerateShortPath()

	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userId in database %v", err)
		return "", err
	}
	err = dbimpl.CreateShortURL(s.data, uuid.String(), shortURL, originalURL)
	if err != nil {
		if errors.Is(err, config.ErrExists) {
			existingShortURL, err := dbimpl.GetShortURLByOriginalURL(s.data, originalURL)
			if err != nil {
				return "", err
			}
			return config.BaseURL + "/" + existingShortURL, nil
		}
		return "", err
	}
	return config.BaseURL + "/" + shortURL, nil
}

// GetOriginalLink retrieves the original URL from the database for a given short URL.
func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := dbimpl.GetOriginalURL(s.data, shortURL)
	if err != nil {
		logger.Errorf("Error retrieving original URL: %v", err)
		return "", err
	}
	return originalURL, nil
}

// Ping checks the connectivity and status of the database.
func (s *service) Ping(ctx context.Context) (int, error) {
	err := s.data.Ping(context.Background())
	if err != nil {
		logger.Errorf("Failed to ping the database: %v", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Close cleans up resources associated with the service, particularly closing any open database connections.
func (s *service) Close() error {
	err := s.data.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetAllURLs retrieves all URLs associated with a specific user ID from the database.
func (s *service) GetAllURLS(ctx context.Context, userID, baseURL string) ([]models.UserURLs, error) {
	res, err := dbimpl.GetOriginalURLsByUserID(s.data, userID, baseURL)
	if err != nil {
		logger.Errorf("Error retrieving all user URLs: %v", err)
		return nil, err
	}
	return res, nil
}

// MarkURLsAsDeleted marks specified URLs as deleted in the database for a given user ID.
func (s *service) MarkURLsAsDeleted(ctx context.Context, userID string, shortURLs []string) error {
	dbimpl.MarkDeleted(s.data, userID, shortURLs)
	return nil
}

// GetStats provides statistics about the service, such as the number of shortened URLs and registered users.
func (s *service) GetStats(ctx context.Context) (urlsCount int, usersCount int, err error) {
	return dbimpl.GetStats(s.data, ctx)
}
