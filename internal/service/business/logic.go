package business

import (
	"math/rand"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
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
	if cache.MockCacheURL != nil {
		return cache.MockCacheURL(originalURL)
	}

	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	shortURL := GenerateShortPath()
	for _, exists := cache.Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}
	cache.Cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL
}

func GetOriginalURL(shortURL string) (string, bool) {
	cache.Mu.RLock()
	defer cache.Mu.RUnlock()

	originalURL, exists := cache.Cache[shortURL]
	return originalURL, exists
}
