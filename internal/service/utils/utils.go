// Package utils provides utility functions for file-based operations related to URL management,
// including saving, loading, and marking URLs as deleted within JSON files.
package utils

import (
	"math/rand"

	"github.com/gleb-korostelev/short-url/internal/config"
)

// GenerateShortPath generates a random string of a predefined length using a specified
// character set. This string is typically used as a short identifier for URLs.
//
// Returns:
//
//	A string representing the randomly generated path.
//
// The function uses the character set and length specified in the config package to
// generate the path. Each character in the path is randomly selected from the character set.
func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

// CheckURL checks if a specific URL (check) is present in a list of URLs (findlist).
//
// Parameters:
//
//	check: The URL to find in the list.
//	findlist: A list of URLs to search through.
//
// Returns:
//
//	True if the URL is found in the list; otherwise, false.
//
// This function iterates over the list of URLs and returns true if the specified URL is
// found within the list. It is commonly used to check for the existence of a URL in a
// list of URLs that may need to be processed or excluded from processing.
func CheckURL(check string, findlist []string) bool {
	for _, find := range findlist {
		if find == check {
			return true
		}
	}
	return false
}
