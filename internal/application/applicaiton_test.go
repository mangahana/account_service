package application

import (
	"account/internal/application/dtos"
	"account/internal/application/mock"
	"account/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setup(t *testing.T) (
	context.Context, *zap.Logger,
	*mock.MockUserService, *mock.MockBanService,
	*mock.MockCodeService, *mock.MockSessionService,
	*mock.MockStorageService,
) {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	// mocks
	ctrl := gomock.NewController(t)

	userService := mock.NewMockUserService(ctrl)
	banService := mock.NewMockBanService(ctrl)
	codeService := mock.NewMockCodeService(ctrl)
	sessionService := mock.NewMockSessionService(ctrl)
	storageService := mock.NewMockStorageService(ctrl)

	return ctx, logger, userService, banService, codeService, sessionService, storageService
}

func TestRegister(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		userService.EXPECT().IsPhoneExists(c, "7775559966").Return(false, nil)
		codeService.EXPECT().Send(c, "7775559966", "127.0.0.1").Return(nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.RegisterInput{Phone: "7775559966", IP: "127.0.0.1"}
		err := app.Register(c, dto)

		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		userService.EXPECT().IsPhoneExists(c, "7775559966").Return(true, nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.RegisterInput{Phone: "7775559966"}
		err := app.Register(c, dto)

		assert.ErrorIs(t, err, domain.ErrPhoneAlreadyInUse)
	})
}

func TestConfirmCode(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		codeService.EXPECT().Verify(c, "7778889966", "1234").Return(nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.ConfirmCodeInput{Phone: "7778889966", Code: "1234"}
		err := app.ConfirmCode(c, dto)

		assert.NoError(t, err)
	})
}

func TestCompleteRegister(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		// setup mocks
		codeService.EXPECT().Verify(c, "7778889966", "1234").Return(nil)
		codeService.EXPECT().RemoveAll(c, "7778889966").Return(nil)
		userService.EXPECT().Create(c, gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
		sessionService.EXPECT().Create(c, 1).Return(&dtos.AuthOutput{AccessToken: "this is a token"}, nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

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
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		user, _ := domain.NewUser("john", "7775556699", "12345678")
		userService.EXPECT().FindOneByPhone(c, "7775556699").Return(user, nil)
		sessionService.EXPECT().Create(c, gomock.Any()).Return(&dtos.AuthOutput{AccessToken: "this is a token"}, nil)

		service := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.LoginInput{Phone: "7775556699", Password: "12345678"}
		output, err := service.Login(c, dto)

		assert.NoError(t, err)
		assert.NotZero(t, output)
	})

	t.Run("fail", func(t *testing.T) {
		user, _ := domain.NewUser("john", "7775556699", "12345678")
		userService.EXPECT().FindOneByPhone(c, "7775556699").Return(user, nil)
		service := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.LoginInput{Phone: "7775556699", Password: "wrongpass"}
		_, err := service.Login(c, dto)

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})
}

func TestRecovery(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		userService.EXPECT().IsPhoneExists(c, "7779998866").Return(true, nil)
		codeService.EXPECT().Send(c, "7779998866", "127.0.0.1").Return(nil)
		service := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := dtos.RecoveryInput{
			Phone: "7779998866",
			IP:    "127.0.0.1",
		}

		err := service.Recovery(c, &dto)
		assert.NoError(t, err)
	})
}

func TestCompleteRecovery(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		codeService.EXPECT().Verify(c, "7778889966", "1234").Return(nil)
		userService.EXPECT().FindOneByPhone(c, "7778889966").Return(&domain.User{ID: 1}, nil)
		userService.EXPECT().UpdatePassword(c, 1, gomock.Any()).Return(nil)
		sessionService.EXPECT().Create(c, gomock.Any()).Return(&dtos.AuthOutput{AccessToken: "this is a token"}, nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.CompleteRecovery{Phone: "7778889966", Code: "1234", Password: "qwerty123"}
		session, err := app.CompleteRecovery(c, dto)

		assert.NoError(t, err)
		assert.NotZero(t, session)
	})
}

func TestBan(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		callerUser := &domain.User{Role: domain.Role{ID: 3}}
		targetUser := &domain.User{Role: domain.Role{ID: 1}}

		userService.EXPECT().FindOneByID(c, gomock.Any()).Return(targetUser, nil)
		userService.EXPECT().FindOneByID(c, gomock.Any()).Return(callerUser, nil)

		expiry := time.Now().Add(time.Second * 15)
		banService.EXPECT().Ban(c, gomock.Any(), gomock.Any(), "toxic", expiry).Return(nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.BanInput{
			CallerUserID: 2,
			UserID:       2,
			Reason:       "toxic",
			Expiry:       expiry,
		}
		err := app.Ban(c, dto)

		assert.NoError(t, err)
	})
}

func TestUnBan(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		banService.EXPECT().FindOneByID(c, 1).Return(&domain.Ban{ID: 1, BannedByID: 2}, nil)
		banService.EXPECT().UnBan(c, 1, 1, "mistake").Return(nil)

		callerUser := &domain.User{Role: domain.Role{ID: 3}}
		userService.EXPECT().FindOneByID(c, 1).Return(callerUser, nil)

		banCallerUser := &domain.User{Role: domain.Role{ID: 2}}
		userService.EXPECT().FindOneByID(c, 2).Return(banCallerUser, nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)

		dto := &dtos.UnBanInput{
			UserID: 1,
			BanID:  1,
			Reason: "mistake",
		}
		err := app.UnBan(c, dto)

		assert.NoError(t, err)
	})
}

func TestAuthenticate(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		user := &domain.User{
			ID:   1,
			Role: domain.Role{Permissions: []string{}},
		}
		userService.EXPECT().FindOneByAccessToken(c, "token").Return(user, nil)
		banService.EXPECT().IsUserBanned(c, 1).Return(false, nil)

		app := New(logger, userService, banService, codeService, sessionService, storageService)
		output, err := app.Authenticate(c, &dtos.AuthenticateInput{AccessToken: "token"})

		assert.NoError(t, err)
		assert.NotZero(t, output)
	})
}

func TestFindByID(t *testing.T) {
	c, logger, userService, banService, codeService, sessionService, storageService := setup(t)

	t.Run("success", func(t *testing.T) {
		returnValue, err := domain.NewUser("john", "777777777", "qwerty123")
		if err != nil {
			t.Fatal(err)
		}

		userService.EXPECT().FindOneByID(c, 1).Return(returnValue, nil)
		app := New(logger, userService, banService, codeService, sessionService, storageService)

		user, err := app.FindByID(c, &dtos.FindByIDInput{ID: 1})

		assert.NoError(t, err)
		assert.NotZero(t, user)
	})
}
