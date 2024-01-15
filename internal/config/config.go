package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string      `yaml:"env" env-default:"local"`
	Grpc GRPC_server `yaml:"grpc_server"`
	Http HTTPServer  `yaml:"http_server"`
}

type GRPC_server struct {
	Address   string        `yaml:"address" env-default:"localhost:4545"`
	Timeout   time.Duration `yaml:"timeout" env-default:"10s"`
	Max_retry int           `yaml:"max_retry" env-default:"5"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func ReadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_CLIENT")

	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_CLIENT is wrong")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ConfigClient file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("can not read config %w", err)
	}

	return &cfg, nil
}
