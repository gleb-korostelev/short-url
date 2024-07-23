// Package cache provides in-memory storage mechanisms for URL shortening and management.
package cache

import "github.com/gleb-korostelev/short-url.git/internal/models"

// Cache is an in-memory storage map that holds the short to original URL mappings.
// The map keys are the shortened URL strings, and the values are URLData instances
// which contain detailed information about the URLs.
var (
	Cache = make(map[string]models.URLData)
)
