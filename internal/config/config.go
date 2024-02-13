package config

import (
	"flag"
	"sync"
)

const (
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length  = 8
)

var (
	ServerAddr string
	BaseURL    string
)

func init() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "address to run HTTP server on")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "base address for the resulting shortened URLs")
}

var (
	Cache = make(map[string]string)
	Mu    sync.RWMutex
)
