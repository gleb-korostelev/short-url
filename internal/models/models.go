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
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

type ShortenBatchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type AllUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
