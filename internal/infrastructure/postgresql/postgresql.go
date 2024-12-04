package postgresql

import (
	"account/internal/infrastructure/configuration"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *configuration.DBConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Name)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
