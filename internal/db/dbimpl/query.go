// Package dbimpl contains the implementation of the database interface defined in the db package.
// It includes functions for initializing the database tables, creating, retrieving, and updating URL data.
package dbimpl

import (
	"context"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

// InitializeTables creates the necessary database tables if they do not already exist.
// This function is typically called at application startup.
func InitializeTables(db db.DB) error {
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

// CreateShortURL inserts a new shortened URL into the database.
// It handles conflicts by updating existing entries where the original URL is already present but marked as deleted.
func CreateShortURL(db db.DB, uuid, shortURL, originalURL string) error {
	sql := `
    INSERT INTO shortened_urls (user_id, short_url, original_url, is_deleted)
    VALUES ($1, $2, $3, FALSE)
    ON CONFLICT (original_url)
    DO UPDATE SET 
        user_id = EXCLUDED.user_id,
        is_deleted = FALSE
    WHERE shortened_urls.is_deleted = TRUE
`
	cmdTag, err := db.Exec(context.Background(), sql, uuid, shortURL, originalURL)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return config.ErrExists
	}
	return nil
}

// GetOriginalURL retrieves the original URL from a shortened URL.
// It returns an error if the URL is marked as deleted or if the shortened URL does not exist.
func GetOriginalURL(db db.DB, shortURL string) (string, error) {
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

// GetOriginalURLsByUserID retrieves all active (not deleted) original URLs for a given user ID.
// It prepends the base URL to each short URL before returning the list.
func GetOriginalURLsByUserID(db db.DB, userID, baseURL string) ([]models.UserURLs, error) {
	sql := `
	SELECT short_url, original_url FROM shortened_urls
	WHERE user_id = $1 AND is_deleted = FALSE
	`
	rows, err := db.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.UserURLs
	for rows.Next() {
		var data models.UserURLs
		if err := rows.Scan(&data.ShortURL, &data.OriginalURL); err != nil {
			return nil, err
		}
		data.ShortURL = baseURL + "/" + data.ShortURL
		urls = append(urls, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

// GetShortURLByOriginalURL retrieves the shortened URL for a given original URL.
// It returns an error if the original URL does not exist in the database.
func GetShortURLByOriginalURL(db db.DB, originalURL string) (string, error) {
	var shortURL string
	sql := `SELECT short_url FROM shortened_urls WHERE original_url = $1`
	err := db.QueryRow(context.Background(), sql, originalURL).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

// MarkDeleted marks a list of shortened URLs as deleted for a specific user.
// This function runs asynchronously and logs the result of the operation.
func MarkDeleted(db db.DB, userID string, shortURLs []string) {
	go func() {
		sql := `
		UPDATE shortened_urls SET is_deleted = TRUE
		WHERE user_id = $1 AND short_url = ANY($2)
		`
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
