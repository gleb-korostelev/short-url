package storage

import (
	"context"
)

type Storage interface {
	SaveUniqueURL(ctx context.Context, originalURL string) (string, int, error)
	SaveURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalLink(ctx context.Context, shortURL string) (string, error)
	Ping(ctx context.Context) (int, error)
	Close() error
}
