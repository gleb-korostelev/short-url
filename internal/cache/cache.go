package cache

import "sync"

var (
	Cache = make(map[string]string)
	Mu    sync.RWMutex
)
