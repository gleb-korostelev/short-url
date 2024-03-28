package dbimpl

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
)

func InitializeTables(db db.DatabaseI) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS shortened_urls (
        id SERIAL PRIMARY KEY,
		uuid UUID NOT NULL UNIQUE,
        short_url VARCHAR(255) UNIQUE NOT NULL,
        original_url VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(context.Background(), createTableSQL)
	return err
}

func CreateShortURL(db db.DatabaseI, uuid, shortURL, originalURL string) error {
	sql := `INSERT INTO shortened_urls (uuid, short_url, original_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING`
	cmdTag, err := db.Exec(context.Background(), sql, uuid, shortURL, originalURL)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return config.ErrExists
	}
	return nil
}

func GetOriginalURL(db db.DatabaseI, shortURL string) (string, error) {
	var originalURL string
	sql := `SELECT original_url FROM shortened_urls WHERE short_url = $1`
	err := db.QueryRow(context.Background(), sql, shortURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func GetOriginalURLByUUID(db db.DatabaseI, uuid string) (string, error) {
	var originalURL string
	sql := `SELECT original_url FROM shortened_urls WHERE uuid = $1`
	err := db.QueryRow(context.Background(), sql, uuid).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func GetShortURLByOriginalURL(db db.DatabaseI, originalURL string) (string, error) {
	var shortURL string
	sql := `SELECT short_url FROM shortened_urls WHERE original_url = $1`
	err := db.QueryRow(context.Background(), sql, originalURL).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}
