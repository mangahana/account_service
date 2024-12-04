package code

import (
	"account/internal/domain"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) FindLatestByIP(c context.Context, ip string, timestamp time.Time) ([]domain.Code, error) {
	output := []domain.Code{}

	sql := "SELECT code, phone, ip, created_at FROM codes WHERE ip = $1 AND created_at > $2;"
	rows, err := r.db.Query(c, sql, ip, timestamp.UTC())
	if err != nil {
		return output, nil
	}

	for rows.Next() {
		var c domain.Code
		if err := rows.Scan(&c.Code, &c.Phone, &c.IP, &c.CreatedAt); err != nil {
			return output, nil
		}
		output = append(output, c)
	}

	return output, nil
}

func (r *repo) FindOneByCredentials(c context.Context, phone, code string) (*domain.Code, error) {
	var output domain.Code

	sql := "SELECT code, phone, ip, created_at FROM codes WHERE phone = $1 AND code = $2;"
	err := r.db.QueryRow(c, sql, phone, code).Scan(&output.Code, &output.Phone, &output.IP, &output.CreatedAt)

	return &output, err
}

func (r *repo) FindOneByPhoneAndIP(c context.Context, phone, ip string) (*domain.Code, error) {
	var output domain.Code

	sql := "SELECT code, phone, ip, created_at FROM codes WHERE phone = $1 AND ip = $2;"
	err := r.db.QueryRow(c, sql, phone, ip).Scan(&output.Code, &output.Phone, &output.IP, &output.CreatedAt)

	return &output, err
}

func (r *repo) Save(c context.Context, code *domain.Code) error {
	sql := "INSERT INTO codes (code, phone, ip) VALUES ($1, $2, $3);"
	_, err := r.db.Exec(c, sql, code.Code, code.Phone, code.IP)
	return err
}

func (r *repo) RemoveAll(c context.Context, phone string) error {
	sql := "DELETE FROM codes WHERE phone = $1;"
	_, err := r.db.Exec(c, sql, phone)
	return err
}
