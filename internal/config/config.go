package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT"`

	JWTSecret string `env:"JWT_SECRET"`

	PostgresHost string `env:"POSTGRES_HOST"`
	PostgresPort string `env:"POSTGRES_PORT"`
	PostgresDB   string `env:"POSTGRES_DB"`
	PostgresUser string `env:"POSTGRES_USER"`
	PostgresPass string `env:"POSTGRES_PASSWORD"`

	AccessTTL  time.Duration `env:"ACCESS_TOKEN_TTL"`
	RefreshTTL time.Duration `env:"REFRESH_TOKEN_TTL"`
}

func MustLoad() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		panic(err)
	}

	return cfg
}