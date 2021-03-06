package milight

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port   int    `json:"port"`
	Bridge string `json:"bridge"`
}

func NewConfig(configPath string) (Config, error) {
	c := Config{}
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return c, err
	}
	file, _ := os.Open(absConfigPath)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}
