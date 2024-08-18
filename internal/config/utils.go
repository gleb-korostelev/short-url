package config

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/gleb-korostelev/short-url.git/internal/models"
)

// This function loads settings from JSON file
func LoadConfig(path string) (*models.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data strings.Builder

	for scanner.Scan() {
		data.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var cfg models.Config
	if err := json.Unmarshal([]byte(data.String()), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
