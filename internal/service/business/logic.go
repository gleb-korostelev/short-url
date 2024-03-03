package business

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
)

func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

func CacheURL(originalURL string) (string, error) {
	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	shortURL := GenerateShortPath()
	for _, exists := cache.Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}
	cache.Cache[shortURL] = originalURL

	var save models.URLData
	save.OriginalURL = originalURL
	save.ShortURL = shortURL
	save.UUID = fmt.Sprint(len(cache.Cache))

	err := SaveURLs(save)
	if err != nil {
		return "", err
	}

	return config.BaseURL + "/" + shortURL, nil
}

func GetOriginalURL(shortURL string) (string, bool) {
	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	originalURL, exists := cache.Cache[shortURL]
	return originalURL, exists
}

func SaveURLs(save models.URLData) error {
	if config.BaseFilePath == "" {
		return nil
	}
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
	if config.BaseFilePath == "" {
		return nil
	}
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

	return nil
}
