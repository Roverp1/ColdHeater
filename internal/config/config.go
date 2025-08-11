package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		AgingPeriod int `yaml:"aging_period"`
	}
}

func LoadConfig() (*Config, error) {
	var config Config

	configFile, err := os.ReadFile("configs/app.yaml")
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file:\n%v", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse YAML config:\n%v", err)
	}

	return &config, nil
}
