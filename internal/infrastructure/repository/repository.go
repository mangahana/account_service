package repository

import "github.com/jackc/pgx/v5/pgxpool"

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}
