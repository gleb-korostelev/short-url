// Package config manages the configuration settings for the application.
// It supports reading from command-line flags and environment variables to set up server parameters,
// database connection strings, file paths, and other necessary settings.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// List of constants
const (
	// Letters defines the character set for generating random strings (used in URLs).
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Length is the default length for generated strings.
	Length = 8

	// DefaultServerAddress specifies the default address for the HTTP server.
	DefaultServerAddress = "localhost:8080"

	// DefaultBaseURL is the default base URL for shortened URLs.
	DefaultBaseURL = "http://localhost:8080"

	// DefaultFilePath is the default file path for storing URL data in JSON format.
	DefaultFilePath = "./tmp/short-url-db.json"

	// TokenExpirationInHour sets the default expiration time for authentication tokens.
	TokenExpirationInHour = 24

	// MaxConcurrentUpdates defines the maximum number of concurrent update operations.
	MaxConcurrentUpdates = 10

	//Certificate file path
	CertFilePath = "./internal/certs/server.crt"

	// Certificate key file path
	KeyFilePath = "./internal/certs/server.key"
)

type contextKey string

// UserContextKey is a key used for storing user ID in the request context.
const UserContextKey = contextKey("user")

// List of vars
var (
	// ErrExists indicates an error when a URL already exists in the storage.
	ErrExists = errors.New("URL already exists")

	// ErrNotFound indicates an error when a URL does not exist in the storage.
	ErrNotFound = errors.New("URL doesn't exists")

	// ErrWrongMode indicates an error when an invalid mode is used (not database mode).
	ErrWrongMode = errors.New("wrong, non db mode")

	// ErrTokenInvalid indicates an error when the provided authentication token is invalid.
	ErrTokenInvalid = errors.New("token is not valid")

	// ErrGone indicates an error when a link has been marked as deleted.
	ErrGone = errors.New("this link is gone")
)

// Configuration variables are settable via command-line flags or environment variables.
var (
	ServerAddr   string                   // ServerAddr is the address where the HTTP server will run.
	BaseURL      string                   // BaseURL is the base address for resulting shortened URLs.
	BaseFilePath string                   // BaseFilePath is the file path where URLs are stored when file mode is used.
	DBDSN        string                   // DBDSN is the Data Source Name for the database connection.
	JwtKeySecret = "very-very-secret-key" // JwtKeySecret is the secret key for signing JWTs.
	EnableHTTPS  bool                     //EnableHTTPS flag
)

// ConfigInit initializes the application's configuration by parsing command-line flags
// and reading from environment variables. It provides default values and overrides them
// with any user-specified options.
func ConfigInit() {
	flag.StringVar(&ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&BaseURL, "b", DefaultBaseURL, "base address for the resulting shortened URLs")
	flag.StringVar(&BaseFilePath, "f", DefaultFilePath, "base file path to save URLs")
	flag.StringVar(&DBDSN, "d", "", "database connection string")
	flag.BoolVar(&EnableHTTPS, "s", false, "Enable HTTPS")

	flag.Parse()

	// Override default values with environment variables if they exist.
	ServerAddr = GetEnv("SERVER_ADDRESS", ServerAddr)
	BaseURL = GetEnv("BASE_URL", BaseURL)
	BaseFilePath = GetEnv("FILE_STORAGE_PATH", BaseFilePath)
	DBDSN = GetEnv("DATABASE_DSN", DBDSN)
	if os.Getenv("ENABLE_HTTPS") == "true" {
		EnableHTTPS = true
	}
}

// GetEnv retrieves the value of an environment variable or returns a fallback value if not set.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Println(key, "set to", value) // Log the environment variable being used.
		return value
	}
	return fallback
}
