package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	FilePath    string   `yaml:"file_path"`
	ServerPort  string   `yaml:"server_port"`
	RecordTypes []string `yaml:"record_types"`
}

func New(configPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to read config: %s", err.Error())
	}

	return &cfg, nil
}
