package dbimpl

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

func InitializeTables(db db.DatabaseI) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS shortened_urls (
        id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL,
        short_url VARCHAR(255) UNIQUE NOT NULL,
        original_url VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		is_deleted BOOLEAN DEFAULT FALSE
    );`
	_, err := db.Exec(context.Background(), createTableSQL)

	return err
}

func CreateShortURL(db db.DatabaseI, uuid, shortURL, originalURL string) error {
	isDeleted := false
	sql := `INSERT INTO shortened_urls (user_id, short_url, original_url, is_deleted) VALUES ($1, $2, $3, $4) ON CONFLICT (original_url) DO NOTHING`
	cmdTag, err := db.Exec(context.Background(), sql, uuid, shortURL, originalURL, isDeleted)
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
	var isDeleted bool
	sql := `SELECT original_url, is_deleted FROM shortened_urls WHERE short_url = $1`
	err := db.QueryRow(context.Background(), sql, shortURL).Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", err
	}
	if isDeleted {
		return "", config.ErrGone
	}
	return originalURL, nil
}

func GetOriginalURLByUUID(db db.DatabaseI, uuid string) ([]models.AllUserURL, error) {
	sql := `SELECT short_url, original_url FROM shortened_urls WHERE user_id = $1 AND is_deleted = FALSE`
	rows, err := db.Query(context.Background(), sql, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.AllUserURL
	for rows.Next() {
		var data models.AllUserURL
		if err := rows.Scan(&data.ShortURL, &data.OriginalURL); err != nil {
			return nil, err
		}
		data.ShortURL = config.BaseURL + "/" + data.ShortURL
		urls = append(urls, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
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

func MarkDeleted(db db.DatabaseI, userID string, shortURLs []string) {
	go func() {
		sql := "UPDATE shortened_urls SET is_deleted = TRUE WHERE user_id = $1 AND short_url = ANY($2)"
		cmdTag, err := db.Exec(context.Background(), sql, userID, shortURLs)
		if err != nil {
			logger.Errorf("Error marking URLs as deleted: %v\n", err)
			return
		}
		if cmdTag.RowsAffected() == 0 {
			logger.Info("No URLs were marked as deleted.")
		} else {
			logger.Infof("%d URLs were marked as deleted.\n", cmdTag.RowsAffected())
		}
	}()
}
