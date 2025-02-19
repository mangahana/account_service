package application

import (
	"account/internal/application/dtos"
	"account/internal/application/mock"
	"account/internal/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setup(t *testing.T) (context.Context, *zap.Logger, *mock.MockUserService, *mock.MockCodeService, *mock.MockSessionService) {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	// mocks
	ctrl := gomock.NewController(t)

	userService := mock.NewMockUserService(ctrl)
	codeService := mock.NewMockCodeService(ctrl)
	sessionService := mock.NewMockSessionService(ctrl)

	return ctx, logger, userService, codeService, sessionService
}

func TestRegister(t *testing.T) {
	c, logger, userService, codeService, sessionService := setup(t)

	t.Run("success", func(t *testing.T) {
		userService.EXPECT().IsPhoneExists(c, "7775559966").Return(false, nil)
		codeService.EXPECT().Send(c, "7775559966", "127.0.0.1").Return(nil)

		app := New(logger, userService, codeService, sessionService)

		dto := &dtos.RegisterInput{Phone: "7775559966", IP: "127.0.0.1"}
		err := app.Register(c, dto)

		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		userService.EXPECT().IsPhoneExists(c, "7775559966").Return(true, nil)

		app := New(logger, userService, codeService, sessionService)

		dto := &dtos.RegisterInput{Phone: "7775559966"}
		err := app.Register(c, dto)

		assert.ErrorIs(t, err, domain.ErrPhoneAlreadyInUse)
	})
}

func TestConfirmCode(t *testing.T) {
	c, logger, userService, codeService, sessionService := setup(t)

	t.Run("success", func(t *testing.T) {
		codeService.EXPECT().Verify(c, "7778889966", "1234").Return(nil)

		app := New(logger, userService, codeService, sessionService)

		dto := &dtos.ConfirmCodeInput{Phone: "7778889966", Code: "1234"}
		err := app.ConfirmCode(c, dto)

		assert.NoError(t, err)
	})
}

func TestCompleteRegister(t *testing.T) {
	c, logger, userService, codeService, sessionService := setup(t)

	t.Run("success", func(t *testing.T) {
		// setup mocks
		codeService.EXPECT().Verify(c, "7778889966", "1234").Return(nil)
		codeService.EXPECT().RemoveAll(c, "7778889966").Return(nil)
		userService.EXPECT().Create(c, gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
		sessionService.EXPECT().Create(c, 1).Return(&dtos.AuthOutput{AccessToken: "this is a token"}, nil)

		app := New(logger, userService, codeService, sessionService)

		dto := &dtos.CompleteRegisterInput{
			Phone:    "7778889966",
			Code:     "1234",
			Username: "john",
			Password: "12345678",
		}

		authRes, err := app.CompleteRegister(c, dto)

		assert.NoError(t, err)
		assert.NotZero(t, authRes)
	})
}

func TestLogin(t *testing.T) {
	c, logger, userService, codeService, sessionService := setup(t)

	t.Run("success", func(t *testing.T) {
		user, _ := domain.NewUser("john", "7775556699", "12345678")
		userService.EXPECT().FindOneByPhone(c, "7775556699").Return(user, nil)
		sessionService.EXPECT().Create(c, gomock.Any()).Return(&dtos.AuthOutput{AccessToken: "this is a token"}, nil)

		service := New(logger, userService, codeService, sessionService)

		dto := &dtos.LoginInput{Phone: "7775556699", Password: "12345678"}
		output, err := service.Login(c, dto)

		assert.NoError(t, err)
		assert.NotZero(t, output)
	})

	t.Run("fail", func(t *testing.T) {
		user, _ := domain.NewUser("john", "7775556699", "12345678")
		userService.EXPECT().FindOneByPhone(c, "7775556699").Return(user, nil)
		service := New(logger, userService, codeService, sessionService)

		dto := &dtos.LoginInput{Phone: "7775556699", Password: "wrongpass"}
		_, err := service.Login(c, dto)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})
}
