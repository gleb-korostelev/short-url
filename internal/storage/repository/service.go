package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/db/dbimpl"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

type service struct {
	data db.DatabaseI
}

func NewDBStorage(data db.DatabaseI) storage.Storage {
	return &service{
		data: data,
	}
}

func (s *service) SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error) {
	shortURL := business.GenerateShortPath()

	uuid, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("Error with parsing userId in database %v", err)
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

func (s *service) SaveURL(ctx context.Context, originalURL string, userID string) (string, error) {
	shortURL := business.GenerateShortPath()

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

func (s *service) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := dbimpl.GetOriginalURL(s.data, shortURL)
	if err != nil {
		logger.Errorf("Error in getting original URL from database %v", err)
		return "", err
	}
	return originalURL, nil
}

func (s *service) Ping(ctx context.Context) (int, error) {
	err := s.data.Ping(context.Background())
	if err != nil {
		logger.Errorf("Failed to connect to the database %v", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *service) Close() error {
	err := s.data.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllURLS(ctx context.Context, userID string) ([]models.AllUserURL, error) {
	res, err := dbimpl.GetOriginalURLByUUID(s.data, userID)
	if err != nil {
		logger.Errorf("Failed to get all user URLS %v", err)
		return nil, err
	}
	return res, nil
}
