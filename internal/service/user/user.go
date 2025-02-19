package user

import (
	"account/internal/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source ./user.go -destination ./mock/mock.go -package mock
type Repository interface {
	FindOneByPhone(c context.Context, phone string) (*domain.User, error)
	FindOneByUsername(c context.Context, username string) (*domain.User, error)

	Create(c context.Context, user *domain.User) (int, error)
}

type service struct {
	repo Repository
}

func New(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) IsPhoneExists(c context.Context, phone string) (bool, error) {
	_, err := s.repo.FindOneByPhone(c, phone)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *service) Create(c context.Context, username, phone, password string) (int, error) {
	_, err := s.repo.FindOneByUsername(c, username)
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, domain.ErrUsernameAlreadyInUse
	}

	user, err := domain.NewUser(username, phone, password)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(c, user)
}

func (s *service) FindOneByPhone(c context.Context, phone string) (*domain.User, error) {
	user, err := s.repo.FindOneByPhone(c, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
