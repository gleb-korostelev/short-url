package business

import (
	"bufio"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/db/impl"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/google/uuid"
)

func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

func CacheURL(w http.ResponseWriter, originalURL string, data db.DatabaseI) (string, error) {
	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	shortURL := GenerateShortPath()
	for _, exists := cache.Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = uuid.New()

	if config.DBDSN != "" {
		err := impl.CreateShortURL(data, save.UUID.String(), save.ShortURL, save.OriginalURL)
		if err != nil {
			if errors.Is(err, config.ErrExists) {
				existingShortURL, err := impl.GetShortURLByOriginalURL(data, save.OriginalURL)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return "", err
				}
				w.WriteHeader(http.StatusConflict)
				return config.BaseURL + "/" + existingShortURL, err
			}
			w.WriteHeader(http.StatusInternalServerError)
			return "", err
		}
	} else if config.BaseFilePath != "" {
		err := SaveURLs(save)
		if err != nil {
			http.Error(w, "Error with saving", http.StatusBadRequest)
			return "", err
		}
	}
	cache.Cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL, nil
}

func OldCacheURL(originalURL string, data db.DatabaseI) (string, error) {
	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	shortURL := GenerateShortPath()
	for _, exists := cache.Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = uuid.New()

	if config.DBDSN != "" {
		err := impl.CreateShortURL(data, save.UUID.String(), save.ShortURL, save.OriginalURL)
		if err != nil {
			logger.Errorf("Error with saving in database %v", err)
			return "", err
		}
	} else if config.BaseFilePath != "" {
		err := SaveURLs(save)
		if err != nil {
			logger.Errorf("Error with saving in file %v", err)
			return "", err
		}
	}
	cache.Cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL, nil
}

func GetOriginalURL(data db.DatabaseI, shortURL string) (string, bool) {
	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	if config.DBDSN != "" {
		originalURL, err := impl.GetOriginalURL(data, shortURL)
		if err != nil {
			logger.Errorf("Error in getting original URL from database %v", err)
			return "", false
		}
		return originalURL, true
	}
	originalURL, exists := cache.Cache[shortURL]
	return originalURL, exists
}

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

func LoadURLs() error {
	if config.BaseFilePath != "" && config.DBDSN == "" {
		file, err := os.Open(config.BaseFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var urlData models.URLData
			if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
				return err
			}
			cache.Cache[urlData.ShortURL] = urlData.OriginalURL
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}
