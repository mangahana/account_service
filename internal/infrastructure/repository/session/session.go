package session

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) Create(c context.Context, userId int, accessToken string) error {
	sql := "INSERT INTO sessions (user_id, access_token) VALUES($1, $2);"
	_, err := r.db.Exec(c, sql, userId, accessToken)
	return err
}
