package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	Letters              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length               = 8
	DefaultServerAddress = "localhost:8081"
	DefaultBaseURL       = "http://localhost:8081"
)

var (
	ServerAddr string
	BaseURL    string
)

func ConfigInit() {

	flag.StringVar(&ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&BaseURL, "b", DefaultBaseURL, "base address for the resulting shortened URLs")

	flag.Parse()

	ServerAddr = GetEnv("SERVER_ADDRESS", ServerAddr)
	BaseURL = GetEnv("BASE_URL", BaseURL)
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(key)
		return value
	}
	return fallback
}
