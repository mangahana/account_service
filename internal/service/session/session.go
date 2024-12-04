package session

import (
	"account/internal/application/dtos"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

//go:generate mockgen -source ./session.go -destination ./mock/mock.go -package mock
type Repository interface {
	Create(c context.Context, userId int, accessToken string) error
}

type service struct {
	repo Repository
}

func New(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(c context.Context, userId int) (*dtos.AuthOutput, error) {
	var output dtos.AuthOutput

	accessToken, err := s.generateRandomToken()
	if err != nil {
		return nil, err
	}

	output.AccessToken = accessToken

	err = s.repo.Create(c, userId, accessToken)
	return &output, err
}

func (s *service) generateRandomToken() (string, error) {
	randomData := make([]byte, 256)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(randomData)
	return hex.EncodeToString(hash[:]), nil
}
