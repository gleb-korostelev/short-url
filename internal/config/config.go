package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	Letters               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length                = 8
	DefaultServerAddress  = "localhost:8080"
	DefaultBaseURL        = "http://localhost:8080"
	DefaultFilePath       = "./tmp/short-url-db.json"
	TokenExpirationInHour = 24
)

type contextKey string

const UserContextKey = contextKey("user")

var (
	ErrExists       = errors.New("URL already exists")
	ErrNotFound     = errors.New("URL doesn't exists")
	ErrWrongMode    = errors.New("wrong, non db mode")
	ErrTokenInvalid = errors.New("token is not valid")
)

var (
	ServerAddr   string
	BaseURL      string
	BaseFilePath string
	DBDSN        string
	JwtKeySecret = "very-very-secret-key"
)

func ConfigInit() {

	flag.StringVar(&ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&BaseURL, "b", DefaultBaseURL, "base address for the resulting shortened URLs")
	flag.StringVar(&BaseFilePath, "f", DefaultFilePath, "base file path to save URLs")
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
