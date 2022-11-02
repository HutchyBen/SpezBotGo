package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DToken     string `json:"dToken"`
	LLPassword string `json:"llPassword"`
}

func EnvironmentConfig() (*Config, error) {
	var config Config

	config.DToken = os.Getenv("DToken")
	config.LLPassword = os.Getenv("LLPassword")
	if config.DToken == "" || config.LLPassword == "" {
		return nil, fmt.Errorf("DToken and LLPassword must be set")
	}
	return &Config{}, nil
}

func NewConfig(fileName string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(fileName)
	if err != nil {
		return EnvironmentConfig()
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return EnvironmentConfig()
	}
	return &config, nil
}
