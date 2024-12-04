package user

import (
	"account/internal/domain"
	"account/internal/service/user/mock"
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (context.Context, *mock.MockRepository) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	return ctx, repo
}

func TestIsPhoneExists(t *testing.T) {
	ctx, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().FindOneByPhone(ctx, "7775556699").Return(&domain.User{ID: 1}, nil)
		service := New(repo)

		exists, err := service.IsPhoneExists(ctx, "7775556699")

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("fail", func(t *testing.T) {
		repo.EXPECT().FindOneByPhone(ctx, "7775556699").Return(nil, pgx.ErrNoRows)
		service := New(repo)

		exists, err := service.IsPhoneExists(ctx, "7775556699")

		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestCreate(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().FindOneByUsername(c, "john").Return(nil, pgx.ErrNoRows)
		repo.EXPECT().Create(c, gomock.Any()).Return(1, nil)

		service := New(repo)

		userId, err := service.Create(c, "john", "7773336699", "12345678")

		assert.NoError(t, err)
		assert.NotZero(t, userId)
	})
}

func TestFindOneByPhone(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().FindOneByPhone(c, "7775556699").Return(&domain.User{ID: 1}, nil)

		service := New(repo)

		user, err := service.FindOneByPhone(c, "7775556699")

		assert.NoError(t, err)
		assert.NotZero(t, user)
	})

	t.Run("fail", func(t *testing.T) {
		repo.EXPECT().FindOneByPhone(c, "7775556699").Return(&domain.User{}, domain.ErrUserNotFound)

		service := New(repo)

		_, err := service.FindOneByPhone(c, "7775556699")

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}
