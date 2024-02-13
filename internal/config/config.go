package config

import (
	"flag"
	"sync"
)

const (
	Letters              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length               = 8
	DefaultServerAddress = "localhost:8080"
	DefaultBaseURL       = "http://localhost:8080"
)

var (
	ServerAddr string
	BaseURL    string
)

func init() {
	flag.StringVar(&ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&BaseURL, "b", DefaultBaseURL, "base address for the resulting shortened URLs")
}

var (
	Cache = make(map[string]string)
	Mu    sync.RWMutex
)
