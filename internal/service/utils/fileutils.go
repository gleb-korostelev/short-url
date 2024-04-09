package utils

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

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
