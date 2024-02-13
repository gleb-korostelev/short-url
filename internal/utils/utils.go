package utils

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

var MockCacheURL func(string) string

func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

func CacheURL(originalURL string) string {
	if MockCacheURL != nil {
		return MockCacheURL(originalURL)
	}

	config.Mu.Lock()
	defer config.Mu.Unlock()

	shortURL := GenerateShortPath()
	for _, exists := config.Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}
	config.Cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL
}

func GetOriginalURL(shortURL string) (string, bool) {
	config.Mu.RLock()
	defer config.Mu.RUnlock()

	originalURL, exists := config.Cache[shortURL]
	return originalURL, exists
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(key)
		return value
	}
	return fallback
}
