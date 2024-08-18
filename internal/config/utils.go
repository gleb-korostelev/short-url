package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gleb-korostelev/short-url.git/internal/models"
)

func LoadConfig(path string) (*models.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg models.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
