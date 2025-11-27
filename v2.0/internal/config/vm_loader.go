package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadVMConfig(path string) (*VMConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vm config: %w", err)
	}

	var cfg VMConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse vm config: %w", err)
	}

	if cfg.Host == "" || cfg.User == "" || cfg.Port == 0 {
		return nil, fmt.Errorf("invalid vm config: missing host/user/port")
	}

	if cfg.SSHKeyPath == "" && !cfg.UseSSHAgent {
		return nil, fmt.Errorf("vm config invalid: either ssh_key_path or use_ssh_agent must be set")
	}

	return &cfg, nil
}
