package ban_repository

import (
	"account/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) FindOneByID(c context.Context, id int) (*domain.Ban, error) {
	var b domain.Ban

	sql := `SELECT id, banned_by, user_id, reason,
						expiry, is_active, created_at
					FROM bans WHERE id = $1;`

	row := r.db.QueryRow(c, sql, id)
	err := row.Scan(
		&b.ID, &b.BannedByID, &b.UserID,
		&b.Reason, &b.Expiry,
		&b.IsActive, &b.Expiry,
	)
	return &b, err
}

func (r *repo) Save(c context.Context, ban *domain.Ban) error {
	sql := `UPDATE bans SET banned_by = $2, user_id = $3, reason = $4,
						expiry = $5, is_active = $6, created_at = $7
					WHERE id = $1;`

	_, err := r.db.Exec(
		c, sql,
		ban.ID, ban.BannedByID, ban.UserID,
		ban.Reason, ban.Expiry, ban.IsActive, ban.CreatedAt,
	)
	return err
}

func (r *repo) Create(c context.Context, ban *domain.Ban) error {
	sql := `INSERT INTO bans (banned_by, user_id, reason, expiry, is_active, created_at)
					 VALUES($1, $2, $3, $4, $5, $6);`
	_, err := r.db.Exec(
		c, sql,
		ban.BannedByID, ban.UserID, ban.Reason,
		ban.Expiry, ban.IsActive, ban.CreatedAt,
	)
	return err
}

func (r *repo) CreateUnban(c context.Context, banId, userId int, reason string) error {
	sql := "INSERT INTO unbans (ban_id, user_id, reason) VALUES ($1, $2, $3);"
	_, err := r.db.Exec(c, sql, banId, userId, reason)
	return err
}
