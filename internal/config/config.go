package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sprint Sprint `yaml:"sprint"`
}

type Sprint struct {
	Start string `yaml:"start"`
	Days  int    `yaml:"days"`
}

func Load(path string) (Config, error) {
	cfg := Config{Sprint: Sprint{Days: 14}}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config.yaml: %w", err)
	}
	if cfg.Sprint.Days <= 0 {
		cfg.Sprint.Days = 14
	}
	return cfg, nil
}
