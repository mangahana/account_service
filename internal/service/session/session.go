package session

import (
	"account/internal/application/dtos"
	"account/internal/common"
	"context"
)

//go:generate mockgen -source ./session.go -destination ./mock/mock.go -package mock
type Repository interface {
	Create(c context.Context, userId int, accessToken string) error
	DeleteAll(c context.Context, userId int) error
}

type service struct {
	repo Repository
}

func New(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(c context.Context, userId int) (*dtos.AuthOutput, error) {
	var output dtos.AuthOutput

	accessToken, err := common.GenerateRandomHash()
	if err != nil {
		return nil, err
	}

	output.AccessToken = accessToken

	err = s.repo.Create(c, userId, accessToken)
	return &output, err
}

func (r *service) Clear(c context.Context, userId int) error {
	return r.repo.DeleteAll(c, userId)
}
