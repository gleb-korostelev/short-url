package business

import (
	"math/rand"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

var (
	Cache        = make(map[string]string)
	Mu           sync.RWMutex
	MockCacheURL func(string) string
)

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

	Mu.RLock()
	defer Mu.RUnlock()

	shortURL := GenerateShortPath()
	for _, exists := Cache[shortURL]; exists; {
		shortURL = GenerateShortPath()
	}
	Cache[shortURL] = originalURL
	return config.BaseURL + "/" + shortURL
}

func GetOriginalURL(shortURL string) (string, bool) {
	Mu.RLock()
	defer Mu.RUnlock()

	originalURL, exists := Cache[shortURL]
	return originalURL, exists
}
