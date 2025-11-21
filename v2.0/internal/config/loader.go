package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads a YAML config file and returns AppConfig.
func Load(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate required fields
	if cfg.AppName == "" || cfg.RepoURL == "" || cfg.ServiceName == "" {
		return nil, errors.New("invalid config: missing app_name, repo_url or service_name")
	}

	return &cfg, nil
}
