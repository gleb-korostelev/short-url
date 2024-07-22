// Package utils provides utility functions for file-based operations related to URL management,
// including saving, loading, and marking URLs as deleted within JSON files.
package utils

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

// SaveURLs writes a URLData object to a file specified in the configuration.
// It appends each new URLData entry as a new line in JSON format.
//
// Parameters:
//
//	save: The URLData object to save to the file.
//
// Returns:
//
//	An error if the file cannot be opened or the data cannot be written; nil otherwise.
func SaveURLs(save models.URLData) error {
	file, err := os.OpenFile(config.BaseFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	data, err := json.Marshal(save)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}

// LoadURLs retrieves the original URL corresponding to a given short URL from a file.
// It scans through each line of the file, looking for a match that is not marked as deleted.
//
// Parameters:
//
//	path: The path to the file containing the URL data.
//	shortURL: The short URL to search for.
//
// Returns:
//
//	The original URL if found and not deleted, or an error if the URL is not found or an error occurs during file processing.
func LoadURLs(path string, shortURL string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var urlData models.URLData
		if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
			return "", err
		}
		if urlData.ShortURL == shortURL && !urlData.DeletedFlag {
			return urlData.OriginalURL, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", config.ErrNotFound
}

// LoadUserURLs retrieves all URLs associated with a specific user ID from a file.
// It only includes URLs that are not marked as deleted.
//
// Parameters:
//
//	path: The path to the file containing the URL data.
//	userID: The user ID to search for.
//
// Returns:
//
//	A list of UserURLs if found, or an error if an error occurs during file processing.
func LoadUserURLs(path string, userID string) ([]models.UserURLs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var urls []models.UserURLs
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var urlData models.URLData
		var data models.UserURLs
		if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
			return nil, err
		}
		if urlData.UUID.String() == userID && !urlData.DeletedFlag {
			data.OriginalURL = urlData.OriginalURL
			data.ShortURL = config.BaseURL + "/" + urlData.ShortURL
			urls = append(urls, data)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

// MarkURLsAsDeletedInFile marks specific URLs as deleted for a given user ID in a file.
// It rewrites the entire file to update the deleted flags of the specified URLs.
//
// Parameters:
//
//	path: The file path where URL data is stored.
//	userID: The user ID whose URLs are to be marked as deleted.
//	shortURLs: A list of short URLs to be marked as deleted.
//
// Returns:
//
//	An error if the file cannot be processed; nil otherwise.
func MarkURLsAsDeletedInFile(path, userID string, shortURLs []string) error {
	file, err := os.OpenFile(config.BaseFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(file)

	for scanner.Scan() {
		var urlData models.URLData
		if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
			return err
		}
		if urlData.UUID.String() == userID && CheckURL(urlData.ShortURL, shortURLs) {
			urlData.DeletedFlag = true
		}
		data, err := json.Marshal(urlData)
		if err != nil {
			logger.Errorf("error marshalling json: %w", err)
			return err
		}
		_, err = writer.WriteString(string(data) + "\n")
		if err != nil {
			logger.Errorf("error writing file: %w", err)
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}
