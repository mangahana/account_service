package session

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool) {
	c := context.Background()

	db, err := pgxpool.New(c, "postgresql://test:test@127.0.0.1/test")
	if err != nil {
		t.Fatal(err)
	}

	db.Exec(c, "TRUNCATE TABLE sessions;")

	return c, db
}

func TestCreate(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)
		err := repo.Create(c, 1, "random token")
		assert.NoError(t, err)
	})
}
