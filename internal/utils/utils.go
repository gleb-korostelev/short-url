package utils

import (
	"math/rand"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

func CacheURL(originalURL string) string {
	config.Mu.Lock()
	defer config.Mu.Unlock()

	shortURL := GenerateShortPath()
	for _, exists := config.Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}
	config.Cache[shortURL] = originalURL
	return shortURL
}

func GetOriginalURL(shortURL string) (string, bool) {
	config.Mu.RLock()
	defer config.Mu.RUnlock()

	originalURL, exists := config.Cache[shortURL]
	return originalURL, exists
}