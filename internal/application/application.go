package application

import (
	"account/internal/application/dtos"
	"account/internal/domain"
	"context"
	"time"

	"go.uber.org/zap"
)

//go:generate mockgen -source ./application.go -destination ./mock/mock.go -package mock
type UserService interface {
	FindOneByID(c context.Context, ID int) (*domain.User, error)
	FindOneByPhone(c context.Context, phone string) (*domain.User, error)

	IsPhoneExists(c context.Context, phone string) (bool, error)

	Create(c context.Context, username, phone, password string) (int, error)

	UpdatePassword(c context.Context, userId int, password string) error
}

type BanService interface {
	Ban(c context.Context, callerId, userId int, reason string, expiry time.Time) error
	UnBan(c context.Context, banId, unBannedByID int, reason string) error

	FindOneByID(c context.Context, id int) (*domain.Ban, error)
}

type CodeService interface {
	Send(c context.Context, phone, ip string) error
	Verify(c context.Context, phone, code string) error
	RemoveAll(c context.Context, phone string) error
}

type SessionService interface {
	Create(c context.Context, userId int) (*dtos.AuthOutput, error)
}

type app struct {
	userService    UserService
	banService     BanService
	codeService    CodeService
	sessionService SessionService
	logger         *zap.Logger
}

func New(
	logger *zap.Logger,
	userService UserService,
	banService BanService,
	codeService CodeService,
	sessionService SessionService,
) *app {
	return &app{
		logger:         logger,
		userService:    userService,
		banService:     banService,
		codeService:    codeService,
		sessionService: sessionService,
	}
}

func (app *app) Register(c context.Context, dto *dtos.RegisterInput) error {
	exists, err := app.userService.IsPhoneExists(c, dto.Phone)
	if err != nil {
		app.logger.Error("failed to check phone", zap.Error(err))
		return err
	}
	if exists {
		return domain.ErrPhoneAlreadyInUse
	}

	if err := app.codeService.Send(c, dto.Phone, dto.IP); err != nil {
		app.logger.Error("failed to send code", zap.Error(err))
		return err
	}

	return nil
}

func (app *app) ConfirmCode(c context.Context, dto *dtos.ConfirmCodeInput) error {
	err := app.codeService.Verify(c, dto.Phone, dto.Code)
	if err != nil {
		app.logger.Error("failed to confirm code", zap.Error(err))
		return err
	}
	return nil
}

func (app *app) CompleteRegister(c context.Context, dto *dtos.CompleteRegisterInput) (*dtos.AuthOutput, error) {
	err := app.codeService.Verify(c, dto.Phone, dto.Code)
	if err != nil {
		app.logger.Error("failed to confirm code", zap.Error(err))
		return nil, err
	}

	userId, err := app.userService.Create(c, dto.Username, dto.Phone, dto.Password)
	if err != nil {
		app.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	if err := app.codeService.RemoveAll(c, dto.Phone); err != nil {
		app.logger.Error("failed to remove codes after register", zap.Error(err))
	}

	res, err := app.sessionService.Create(c, userId)
	if err != nil {
		app.logger.Error("failed to create session", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (app *app) Login(c context.Context, dto *dtos.LoginInput) (*dtos.AuthOutput, error) {
	user, err := app.userService.FindOneByPhone(c, dto.Phone)
	if err != nil {
		app.logger.Error("failed to find user", zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	if err := user.ComparePassword(dto.Password); err != nil {
		return &dtos.AuthOutput{}, domain.ErrInvalidCredentials
	}

	session, err := app.sessionService.Create(c, user.ID)
	if err != nil {
		app.logger.Error("failed to create access token", zap.Error(err))
		return nil, err
	}

	return session, nil
}

func (app *app) Recovery(c context.Context, dto *dtos.RecoveryInput) error {
	exists, err := app.userService.IsPhoneExists(c, dto.Phone)
	if err != nil {
		app.logger.Error("failed to check phone", zap.Error(err))
		return err
	}

	if !exists {
		return domain.ErrPhoneNotFound
	}

	if err := app.codeService.Send(c, dto.Phone, dto.IP); err != nil {
		app.logger.Error("failed to send recovery code", zap.Error(err))
		return err
	}

	return nil
}

func (app *app) CompleteRecovery(c context.Context, dto *dtos.CompleteRecovery) (*dtos.AuthOutput, error) {
	err := app.codeService.Verify(c, dto.Phone, dto.Code)
	if err != nil {
		app.logger.Error("failed to verify code", zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	user, err := app.userService.FindOneByPhone(c, dto.Phone)
	if err != nil {
		app.logger.Error("failed to find user by phone", zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	if err := app.userService.UpdatePassword(c, user.ID, dto.Password); err != nil {
		app.logger.Error("failed to update password", zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	session, err := app.sessionService.Create(c, user.ID)
	if err != nil {
		app.logger.Error("failed to create session", zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	return session, nil
}

func (app *app) Ban(c context.Context, dto *dtos.BanInput) error {
	targetUser, err := app.userService.FindOneByID(c, dto.UserID)
	if err != nil {
		app.logger.Error("failed to find user by ID", zap.Error(err))
		return err
	}

	callerUser, err := app.userService.FindOneByID(c, dto.CallerUserID)
	if err != nil {
		app.logger.Error("failed to find user by ID", zap.Error(err))
		return err
	}

	if callerUser.Role.ID <= targetUser.Role.ID {
		app.logger.Error("cannot ban higher", zap.Error(domain.ErrCannotBanHigherUser))
		return domain.ErrCannotBanHigherUser
	}

	err = app.banService.Ban(c, dto.CallerUserID, dto.UserID, dto.Reason, dto.Expiry)
	if err != nil {
		app.logger.Error("cannot ban user", zap.Error(err))
		return err
	}

	return nil
}

func (app *app) UnBan(c context.Context, dto *dtos.UnBanInput) error {
	ban, err := app.banService.FindOneByID(c, dto.BanID)
	if err != nil {
		app.logger.Error("failed to find ban by id", zap.Error(err))
		return err
	}

	callerUser, err := app.userService.FindOneByID(c, dto.UserID)
	if err != nil {
		app.logger.Error("failed to find user by id", zap.Error(err))
		return err
	}

	banCallerUser, err := app.userService.FindOneByID(c, ban.BannedByID)
	if err != nil {
		app.logger.Error("failed to find user by id", zap.Error(err))
		return err
	}

	if callerUser.ID != banCallerUser.ID {
		if callerUser.Role.ID <= banCallerUser.Role.ID {
			app.logger.Error("failed to unban", zap.Error(domain.ErrCannotBanHigherUser))
			return domain.ErrCannotBanHigherUser
		}
	}

	err = app.banService.UnBan(c, dto.BanID, dto.UserID, dto.Reason)
	if err != nil {
		app.logger.Error("failed to unban user", zap.Error(err))
		return err
	}

	return nil
}
