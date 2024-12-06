package ban_repository

import (
	"account/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (context.Context, *pgxpool.Pool) {
	c := context.Background()

	db, err := pgxpool.New(c, "postgresql://test:test@127.0.0.1/test")
	if err != nil {
		t.Fatal(err)
	}

	db.Exec(c, "TRUNCATE TABLE bans; TRUNCATE TABLE unbans;")

	return c, db
}

func TestFindOneByID(t *testing.T) {
	c, db := setup(t)
	t.Run("success", func(t *testing.T) {
		var banId int

		db.QueryRow(c, `INSERT INTO bans (banned_by, user_id, reason, expiry) VALUES(1, 2, 'toxic', $1) RETURNING id;`, time.Now()).Scan(&banId)

		repo := New(db)
		ban, err := repo.FindOneByID(c, banId)

		assert.NoError(t, err)
		assert.NotZero(t, ban)
	})
}

func TestCreate(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		ban, err := domain.NewBan(1, 2, "toxic", time.Now().Add(time.Second*15))
		if err != nil {
			t.Fatal(err)
		}

		err = repo.Create(c, ban)

		assert.NoError(t, err)
	})
}

func TestCreateUnban(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		err := repo.CreateUnban(c, 1, 2, "wasd")

		assert.NoError(t, err)
	})
}
