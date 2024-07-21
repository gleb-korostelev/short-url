package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type URLPayload struct {
	URL string `json:"url"`
}

type ShortURLResponse struct {
	Result string `json:"result"`
}

type URLData struct {
	UUID        uuid.UUID `db:"user_id"`
	ShortURL    string    `db:"short_url"`
	OriginalURL string    `db:"original_url"`
	DeletedFlag bool      `db:"is_deleted"`
}

type ShortenBatchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
