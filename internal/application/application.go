package application

import (
	"account/internal/application/dtos"
	"account/internal/common"
	"account/internal/domain"
	"context"
	"net/http"
	"net/url"
	"path"
	"slices"
	"time"

	"go.uber.org/zap"
)

//go:generate mockgen -source ./application.go -destination ./mock/mock.go -package mock
type UserService interface {
	FindOneByID(c context.Context, ID int) (*domain.User, error)
	FindOneByPhone(c context.Context, phone string) (*domain.User, error)
	FindOneByAccessToken(c context.Context, accessToken string) (*domain.User, error)
	FindOneByUsername(c context.Context, username string) (*domain.User, error)

	IsPhoneExists(c context.Context, phone string) (bool, error)

	Create(c context.Context, username, phone, password string) (int, error)

	UpdatePassword(c context.Context, userId int, password string) error
	ChangePassword(c context.Context, userId int, oldPassword, newPassword string) error
	UpdatePhoto(c context.Context, userId int, filename string) (*domain.User, error)
}

type BanService interface {
	Ban(c context.Context, callerId, userId int, reason string, expiry time.Time) error
	UnBan(c context.Context, banId, unBannedByID int, reason string) error

	IsUserBanned(c context.Context, userId int) (bool, error)
	FindOneByID(c context.Context, id int) (*domain.Ban, error)
}

type CodeService interface {
	Send(c context.Context, phone, ip string) error
	Verify(c context.Context, phone, code string) error
	RemoveAll(c context.Context, phone string) error
}

type SessionService interface {
	Create(c context.Context, userId int) (*dtos.AuthOutput, error)
	Clear(c context.Context, userId int) error
}

type StorageService interface {
	Put(c context.Context, data []byte) (filename string, err error)
	Remove(c context.Context, filename string) error
}

var formats = []string{"image/jpeg", "image/png", "image/webp"}

type app struct {
	userService    UserService
	banService     BanService
	codeService    CodeService
	sessionService SessionService
	storageService StorageService
	logger         *zap.Logger
}

func New(
	logger *zap.Logger,
	userService UserService,
	banService BanService,
	codeService CodeService,
	sessionService SessionService,
	storageService StorageService,
) *app {
	return &app{
		logger:         logger,
		userService:    userService,
		banService:     banService,
		codeService:    codeService,
		sessionService: sessionService,
		storageService: storageService,
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

func (app *app) Authenticate(c context.Context, dto *dtos.AuthenticateInput) (*dtos.AuthenticateOutput, error) {
	user, err := app.userService.FindOneByAccessToken(c, dto.AccessToken)
	if err != nil {
		app.logger.Error("failed to find user by access token", zap.Error(err))
		return &dtos.AuthenticateOutput{}, err
	}

	isBanned, err := app.banService.IsUserBanned(c, user.ID)
	if err != nil {
		app.logger.Error("failed to check is_banned", zap.Error(err))
		return &dtos.AuthenticateOutput{}, err
	}

	return &dtos.AuthenticateOutput{
		UserID:   user.ID,
		IsBanned: isBanned,
		Role:     user.Role,
	}, err
}

func (app *app) FindByID(c context.Context, dto *dtos.FindByIDInput) (*dtos.UserOutput, error) {
	user, err := app.userService.FindOneByID(c, dto.ID)
	if err != nil {
		app.logger.Error("failed to find user", zap.Error(err))
		return &dtos.UserOutput{}, err
	}

	return dtos.NewUserOutput(user.ID, user.Username, user.Photo, user.Description, user.Role), nil
}

func (app *app) FindByUsername(c context.Context, dto *dtos.FindByUsernameInput) (*dtos.UserOutput, error) {
	user, err := app.userService.FindOneByUsername(c, dto.Username)
	if err != nil {
		app.logger.Error("failed to find user", zap.Error(err))
		return &dtos.UserOutput{}, err
	}

	return dtos.NewUserOutput(user.ID, user.Username, user.Photo, user.Description, user.Role), nil
}

func (app *app) IsPhoneExists(c context.Context, dto *dtos.IsPhoneExistsInput) error {
	_, err := app.userService.FindOneByPhone(c, dto.Phone)
	if err != nil {
		app.logger.Error("failed to find user", zap.Error(err))
		return err
	}

	return nil
}

func (app *app) ChangePassword(c context.Context, dto *dtos.ChangePasswordInput) (*dtos.AuthOutput, error) {
	if err := app.userService.ChangePassword(c, dto.UserID, dto.OldPassword, dto.NewPassword); err != nil {
		app.logger.Error("failed to change password", zap.Int("user_id", dto.UserID), zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	if !dto.Logout {
		return &dtos.AuthOutput{}, nil
	}

	if err := app.sessionService.Clear(c, dto.UserID); err != nil {
		app.logger.Error("failed to clear session", zap.Int("user_id", dto.UserID), zap.Error(err))
	}

	output, err := app.sessionService.Create(c, dto.UserID)
	if err != nil {
		app.logger.Error("failed to create session", zap.Int("user_id", dto.UserID), zap.Error(err))
		return &dtos.AuthOutput{}, err
	}

	return output, nil
}

func (app *app) UploadImage(c context.Context, dto *dtos.UploadImageInput) (*dtos.UploadImageOutput, error) {
	mime := http.DetectContentType(dto.File)

	if !slices.Contains(formats, mime) {
		return nil, domain.ErrUnsupportedFormat
	}

	resizedFile, err := common.Resize(dto.File, 300)
	if err != nil {
		return nil, err
	}

	filename, err := app.storageService.Put(c, resizedFile)
	if err != nil {
		app.logger.Error("failed to find user", zap.Int("user_id", dto.UserID), zap.Error(err))
		return nil, err
	}

	user, err := app.userService.UpdatePhoto(c, dto.UserID, filename)
	if err != nil {
		app.logger.Error("failed to update photo", zap.Int("user_id", dto.UserID), zap.Error(err))
		return nil, err
	}

	return &dtos.UploadImageOutput{
		Filename: user.Photo,
	}, nil
}

func (app *app) RemoveImage(c context.Context, dto *dtos.RemovePhotoInput) error {
	user, err := app.userService.FindOneByID(c, dto.UserID)
	if err != nil {
		app.logger.Error("failed to find user", zap.Int("user_id", dto.UserID), zap.Error(err))
		return err
	}

	url, err := url.Parse(user.Photo)
	if err != nil {
		app.logger.Error("failed to parse photo path", zap.String("path", user.Photo), zap.Error(err))
		return err
	}

	if err := app.storageService.Remove(c, path.Base(url.Path)); err != nil {
		app.logger.Error("failed to remove file from storage", zap.String("filename", path.Base(url.Path)), zap.Error(err))
		return err
	}

	_, err = app.userService.UpdatePhoto(c, dto.UserID, "")
	if err != nil {
		app.logger.Error("failed to update photo", zap.Int("user_id", dto.UserID), zap.Error(err))
		return err
	}

	return nil
}
