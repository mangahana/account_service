package user

import (
	"account/internal/domain"
	"context"
	"testing"

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

	db.Exec(c, "TRUNCATE TABLE users;")

	return c, db
}

func TestFindOneByPhone(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		_, err := db.Exec(c, "INSERT INTO users (username, phone, password) VALUES('doe', '7775556699', '12345678')")
		if err != nil {
			t.Fatal(err)
		}

		user, err := repo.FindOneByPhone(c, "7775556699")

		assert.NoError(t, err)
		assert.NotZero(t, user)
	})

	t.Run("fail", func(t *testing.T) {
		repo := New(db)
		_, err := repo.FindOneByPhone(c, "wrongnumber")
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestFindOneByUsername(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		_, err := db.Exec(c, "INSERT INTO users (username, phone, password) VALUES('doe', '7775556699', '12345678')")
		if err != nil {
			t.Fatal(err)
		}

		user, err := repo.FindOneByUsername(c, "doe")

		assert.NoError(t, err)
		assert.NotZero(t, user)
	})

	t.Run("fail", func(t *testing.T) {
		repo := New(db)
		_, err := repo.FindOneByUsername(c, "wrongusername")
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func TestFindOneByID(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		_, err := db.Exec(c, "INSERT INTO users (id, username, phone, password) VALUES(99, 'doe', '7775556699', '12345678')")
		if err != nil {
			t.Fatal(err)
		}

		user, err := repo.FindOneByID(c, 99)

		assert.NoError(t, err)
		assert.NotZero(t, user)
	})
}

func TestCreate(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		repo := New(db)

		userId, err := repo.Create(c, &domain.User{Username: "john", Phone: "7776668844", Password: "12345678"})

		assert.NoError(t, err)
		assert.NotZero(t, userId)
	})
}

func TestSave(t *testing.T) {
	c, db := setup(t)

	t.Run("success", func(t *testing.T) {
		db.Exec(c, "INSERT INTO users (id, username, phone, password) VALUES (66,'john', '7779998866', 'oldpass');")

		repo := New(db)

		user, _ := domain.NewUser("newuser", "7778885566", "qwerty156")
		user.ID = 1

		err := repo.Save(c, user)

		assert.NoError(t, err)
	})
}
