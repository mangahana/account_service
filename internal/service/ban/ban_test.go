package ban_service

import (
	"account/internal/domain"
	"account/internal/service/ban/mock"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (context.Context, *mock.MockRepository) {
	c := context.Background()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)

	return c, repo
}

func TestBan(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().Save(c, gomock.Any()).Return(nil)

		service := New(repo)

		err := service.Ban(c, 1, 2, "toxic", time.Now().Add(time.Minute*15))

		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		service := New(repo)

		err := service.Ban(c, 1, 2, "", time.Now().Add(time.Minute*15))

		assert.ErrorIs(t, err, domain.ErrReasonCantBeEmpty)
	})
}

func TestUnBan(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().FindOneByID(c, 1).Return(&domain.Ban{ID: 1}, nil)
		repo.EXPECT().Save(c, gomock.Any()).Return(nil)
		repo.EXPECT().CreateUnban(c, 1, 1, "toxic").Return(nil)

		service := New(repo)

		err := service.UnBan(c, 1, 1, "toxic")

		assert.NoError(t, err)
	})
}

func TestFindOneByID(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().FindOneByID(c, 1).Return(&domain.Ban{ID: 1}, nil)

		service := New(repo)

		ban, err := service.FindOneByID(c, 1)

		assert.NoError(t, err)
		assert.NotZero(t, ban)
	})
}
