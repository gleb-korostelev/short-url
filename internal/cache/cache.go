package cache

import (
	"sync"
)

var (
	// URLs  []models.URLData
	Cache = make(map[string]string)
	Mu    sync.RWMutex
)
