package session

import (
	"account/internal/service/session/mock"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (context.Context, *mock.MockRepository) {
	c := context.Background()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)

	return c, repo
}

func TestCreate(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().Create(c, 1, gomock.Any()).Return(nil)

		service := New(repo)

		output, err := service.Create(c, 1)

		assert.NoError(t, err)
		assert.NotZero(t, output)

		t.Log(output.AccessToken)
	})
}
