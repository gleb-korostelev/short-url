package cache

import "github.com/gleb-korostelev/short-url/internal/models"

var (
	Cache = make(map[string]models.URLData)
)
