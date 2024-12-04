package code

import (
	"account/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool) {
	c := context.Background()

	db, err := pgxpool.New(c, "postgresql://test:test@127.0.0.1/test")
	if err != nil {
		t.Fatal(err)
	}

	db.Exec(c, "TRUNCATE TABLE codes;")

	return c, db
}

func TestFindLatestByIP(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		// setup
		db.Exec(c, "INSERT INTO codes (code, phone, ip) VALUES('7823', '7779995566', '127.0.0.1');")

		codes, err := repo.FindLatestByIP(c, "127.0.0.1", time.Now().UTC().Add(time.Minute*-30))
		assert.NoError(t, err)
		assert.Len(t, codes, 1)
	})
}

func TestSave(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		code := &domain.Code{
			Code:      "1234",
			Phone:     "8889995566",
			IP:        "127.0.0.1",
			CreatedAt: time.Now().UTC(),
		}
		err := repo.Save(c, code)

		assert.NoError(t, err)
	})
}

func TestFindOneByCredentials(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		db.Exec(c, "INSERT INTO codes (code, phone, ip) VALUES ('1234','7778889966','127.0.0.1')")

		code, err := repo.FindOneByCredentials(c, "7778889966", "1234")

		assert.NoError(t, err)
		assert.NotZero(t, code)
	})

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		_, err := repo.FindOneByCredentials(c, "wrongnumber", "1234")

		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestFindOneByPhoneAndIP(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		db.Exec(c, "INSERT INTO codes (code, phone, ip) VALUES ('1234','7778889966','127.0.0.1')")

		code, err := repo.FindOneByPhoneAndIP(c, "7778889966", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotZero(t, code)
	})
}

func TestRemoveAll(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		err := repo.RemoveAll(c, "7778889966")

		assert.NoError(t, err)
	})
}
