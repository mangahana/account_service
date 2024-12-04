package code

import (
	"account/internal/domain"
	"account/internal/infrastructure/configuration"
	"account/internal/service/code/mock"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (context.Context, *mock.MockRepository) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	return ctx, repo
}

func TestSend(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		// setup mocks
		code := domain.Code{Code: "1212", Phone: "7778885566", IP: "127.0.0.1", CreatedAt: time.Now()}
		repo.EXPECT().FindLatestByIP(c, "127.0.0.1", gomock.Any()).Return([]domain.Code{code}, nil)
		repo.EXPECT().FindOneByPhoneAndIP(c, "7778889966", "127.0.0.1").Return(&domain.Code{}, nil)
		repo.EXPECT().Save(c, gomock.Any()).Return(nil)

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code":0}`))
		}))
		defer mockServer.Close()

		cfg := &configuration.SMSConfig{
			ApiKey:    "sometoken",
			ApiDomain: mockServer.URL,
		}
		service := New(cfg, repo)

		err := service.Send(c, "7778889966", "127.0.0.1")
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		repo.EXPECT().FindOneByPhoneAndIP(c, "7778889966", "127.0.0.1").Return(&domain.Code{CreatedAt: time.Now()}, nil)

		cfg := &configuration.SMSConfig{}
		service := New(cfg, repo)

		err := service.Send(c, "7778889966", "127.0.0.1")
		assert.ErrorIs(t, err, domain.ErrCodeSendingLimit)
	})
}

func TestVerify(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().FindOneByCredentials(c, "7778889955", "1234").Return(&domain.Code{Code: "1234"}, nil)

		service := New(&configuration.SMSConfig{}, repo)

		err := service.Verify(c, "7778889955", "1234")

		assert.NoError(t, err)
	})
}

func TestRemoveAll(t *testing.T) {
	c, repo := setup(t)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().RemoveAll(c, "7778889955").Return(nil)

		service := New(&configuration.SMSConfig{}, repo)

		err := service.RemoveAll(c, "7778889955")

		assert.NoError(t, err)
	})
}
