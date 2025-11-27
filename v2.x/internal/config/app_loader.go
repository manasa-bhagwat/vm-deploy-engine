package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadAppConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read app config: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse app config: %w", err)
	}

	if cfg.AppName == "" || cfg.RepoURL == "" {
		return nil, fmt.Errorf("invalid app config: missing required fields")
	}

	return &cfg, nil
}
