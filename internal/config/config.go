package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// FIX 7: docker-compose sets HTTP_PORT as PORT — added PORT as primary env key with
	// HTTP_PORT as fallback so both local .env and docker-compose work without changes
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`

	JWTSecret string `env:"JWT_SECRET"` // no default — required

	// FIX 7: docker-compose sets POSTGRES_* vars; these names match exactly
	PostgresHost string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDB   string `env:"POSTGRES_DB"   envDefault:"auth_db"`
	PostgresUser string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPass string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`

	AccessTTL  time.Duration `env:"ACCESS_TOKEN_TTL"  envDefault:"15m"`
	RefreshTTL time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`
	BcryptCost int           `env:"BCRYPT_COST"       envDefault:"10"`
}

func MustLoad() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
