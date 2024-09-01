// Package models contains data structures used throughout the application.
// These structures are used for handling URL data and user authentication.
package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// URLPayload defines the structure for receiving URLs in requests.
type URLPayload struct {
	URL string `json:"url"`
}

// ShortURLResponse defines the structure for sending shortened URLs in responses.
type ShortURLResponse struct {
	Result string `json:"result"`
}

// URLData describes the structure of URL data in the database.
type URLData struct {
	UUID        uuid.UUID `db:"user_id"`      // UUID of the user
	ShortURL    string    `db:"short_url"`    // Shortened URL
	OriginalURL string    `db:"original_url"` // Original URL
	DeletedFlag bool      `db:"is_deleted"`   // Flag indicating if the URL is deleted
}

// ShortenBatchRequestItem describes a request item for batch URL shortening.
type ShortenBatchRequestItem struct {
	CorrelationID string `json:"correlation_id"` // Correlation identifier for tracking requests
	OriginalURL   string `json:"original_url"`   // Original URL to be shortened
}

// ShortenBatchResponseItem describes a response item for a batch URL shortening request.
type ShortenBatchResponseItem struct {
	CorrelationID string `json:"correlation_id"` // Correlation identifier from the request
	ShortURL      string `json:"short_url"`      // Shortened URL
}

// UserURLs represents both shortened and original URLs associated with a user.
type UserURLs struct {
	ShortURL    string `json:"short_url"`    // Shortened URL
	OriginalURL string `json:"original_url"` // Original URL
}

// Claims defines custom JWT claims used for authentication.
type Claims struct {
	UserID string `json:"user_id"` // User identifier
	jwt.RegisteredClaims
}

// Config defines server settings thats taken from JSON file
type Config struct {
	ServerAddr    string `json:"server_address"`
	BaseURL       string `json:"base_url"`
	BaseFilePath  string `json:"file_storage_path"`
	DBDSN         string `json:"database_dsn"`
	EnableHTTPS   bool   `json:"enable_https"`
	TrustedSubnet string `json:"trusted_subnet"`
}

// Stats defines fields of statistic for /api/internal/stats handler
type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
