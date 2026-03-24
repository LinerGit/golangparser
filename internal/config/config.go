package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" json:"service_name" required:"true" default:"parser_service"`
	AppEnv      string `env:"APP_ENV" json:"app_environment" required:"true" default:"development"`
	LogLevel    string `env:"LOG_LEVEL" json:"log_level" required:"true" default:"info"`
	DbDsn       string `env:"DB_DSN" json:"db_dsn" required:"true"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using sys env vars")
	}
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
