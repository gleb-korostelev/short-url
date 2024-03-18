package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	Letters              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length               = 8
	DefaultServerAddress = "localhost:8080"
	DefaultBaseURL       = "http://localhost:8080"
	// DefaultFilePath      = "./tmp/short-url-db.json"
)

var (
	ServerAddr   string
	BaseURL      string
	BaseFilePath string
	DBDSN        string
)

func ConfigInit() {

	flag.StringVar(&ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&BaseURL, "b", DefaultBaseURL, "base address for the resulting shortened URLs")
	flag.StringVar(&BaseFilePath, "f", "", "base file path to save URLs")
	flag.StringVar(&DBDSN, "d", "", "base file path to save URLs")

	flag.Parse()

	ServerAddr = GetEnv("SERVER_ADDRESS", ServerAddr)
	BaseURL = GetEnv("BASE_URL", BaseURL)
	BaseFilePath = GetEnv("FILE_STORAGE_PATH", BaseFilePath)
	DBDSN = GetEnv("DATABASE_DSN", DBDSN)
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(key)
		return value
	}
	return fallback
}
