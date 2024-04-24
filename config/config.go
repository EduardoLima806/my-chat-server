package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml: "name" env:"APP_NAME"`
		Version string `env-required:"true" yaml: "version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PG struct {
		DatabaseDriver string `env-required:"true" yaml:"database_driver" env:"DATABASE_DRIVER"`
		Host           string `env-required:"true" yaml:"db_host" env:"DB_HOST" env-default:"localhost"`
		Port           string `env-required:"true" yaml:"db_port" env:"DB_PORT" env-default:"5432"`
		User           string `env-required:"true" yaml:"db_user" env:"DB_USER" env-default:"postgres"`
		Password       string `env-required:"true" yaml:"db_password" env:"DB_PASSWORD"`
		Name           string `env-required:"true" yaml:"db_name" env:"DB_NAME" env-default:"postgres"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
