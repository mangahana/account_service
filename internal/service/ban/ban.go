package ban_service

import (
	"account/internal/domain"
	"context"
	"time"
)

//go:generate mockgen -source ./ban.go -destination ./mock/mock.go -package mock
type Repository interface {
	FindOneByID(c context.Context, id int) (*domain.Ban, error)

	Save(c context.Context, ban *domain.Ban) error

	Create(c context.Context, ban *domain.Ban) error
	CreateUnban(c context.Context, banId, userId int, reason string) error
}

type service struct {
	repo Repository
}

func New(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Ban(c context.Context, callerId, userId int, reason string, expiry time.Time) error {
	ban, err := domain.NewBan(callerId, userId, reason, expiry)
	if err != nil {
		return err
	}

	if err := s.repo.Create(c, ban); err != nil {
		return err
	}

	return nil
}

func (s *service) UnBan(c context.Context, banId, unBannedByID int, reason string) error {
	ban, err := s.repo.FindOneByID(c, banId)
	if err != nil {
		return err
	}

	ban.IsActive = false

	if err := s.repo.Save(c, ban); err != nil {
		return err
	}

	if err := s.repo.CreateUnban(c, banId, unBannedByID, reason); err != nil {
		return err
	}

	return nil
}

func (s *service) FindOneByID(c context.Context, id int) (*domain.Ban, error) {
	ban, err := s.repo.FindOneByID(c, id)
	if err != nil {
		return &domain.Ban{}, domain.ErrBanNotFound
	}

	return ban, nil
}
