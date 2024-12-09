package grpc

import (
	"account/internal/application/dtos"
	"context"
)

type UseCase interface {
	Register(c context.Context, dto *dtos.RegisterInput) error
	ConfirmCode(c context.Context, dto *dtos.ConfirmCodeInput) error
	CompleteRegister(c context.Context, dto *dtos.CompleteRegisterInput) (*dtos.AuthOutput, error)
	Login(c context.Context, dto *dtos.LoginInput) (*dtos.AuthOutput, error)
	CompleteRecovery(c context.Context, dto *dtos.CompleteRecovery) (*dtos.AuthOutput, error)
	Recovery(c context.Context, dto *dtos.RecoveryInput) error
	Ban(c context.Context, dto *dtos.BanInput) error
	UnBan(c context.Context, dto *dtos.UnBanInput) error
	Authenticate(c context.Context, dto *dtos.AuthenticateInput) (*dtos.AuthenticateOutput, error)
	FindByID(c context.Context, dto *dtos.FindByIDInput) (*dtos.UserOutput, error)
	FindByUsername(c context.Context, dto *dtos.FindByUsernameInput) (*dtos.UserOutput, error)
	IsPhoneExists(c context.Context, dto *dtos.IsPhoneExistsInput) error
	ChangePassword(c context.Context, dto *dtos.ChangePasswordInput) (*dtos.AuthOutput, error)
	UploadImage(c context.Context, dto *dtos.UploadImageInput) (*dtos.UploadImageOutput, error)
	RemoveImage(c context.Context, dto *dtos.RemovePhotoInput) error
}
