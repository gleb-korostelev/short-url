package storage

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/models"
)

type Storage interface {
	SaveUniqueURL(ctx context.Context, originalURL string, userID string) (string, int, error)
	SaveURL(ctx context.Context, originalURL string, userID string) (string, error)
	GetOriginalLink(ctx context.Context, shortURL string) (string, error)
	Ping(ctx context.Context) (int, error)
	Close() error
	GetAllURLS(ctx context.Context, userID, baseURL string) ([]models.UserURLs, error)
	MarkURLsAsDeleted(ctx context.Context, userID string, shortURLs []string) error
}
