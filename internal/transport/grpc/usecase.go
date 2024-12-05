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
}
