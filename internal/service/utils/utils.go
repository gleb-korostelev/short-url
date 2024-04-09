package utils

import (
	"math/rand"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

func CheckURL(check string, findlist []string) bool {
	for _, find := range findlist {
		if find == check {
			return true
		}
	}
	return false
}
