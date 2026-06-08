package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/LinerGit/go-auth/internal/config"
)

func New(cfg config.Config) (*pgxpool.Pool, error) {

	// FIX 5: added sslmode=disable — without it pgx attempts TLS negotiation which
	// fails against a plain Postgres container and causes a connection error on startup
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	return pgxpool.New(
		context.Background(),
		dsn,
	)
}
