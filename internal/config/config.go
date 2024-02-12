package config

import "sync"

const (
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length  = 8
)

var (
	Cache = make(map[string]string)
	Mu    sync.RWMutex
)
