package user

import (
	"account/internal/domain"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source ./user.go -destination ./mock/mock.go -package mock
type Repository interface {
	FindOneByPhone(c context.Context, phone string) (*domain.User, error)
	FindOneByUsername(c context.Context, username string) (*domain.User, error)
	FindOneByID(c context.Context, userId int) (*domain.User, error)
	FindOneByAccessToken(c context.Context, accessToken string) (*domain.User, error)

	Create(c context.Context, user *domain.User) (int, error)
	Save(c context.Context, user *domain.User) error
}

type service struct {
	repo       Repository
	cdnBaseUrl string
}

func New(repo Repository, cdnBaseUrl string) *service {
	return &service{repo: repo, cdnBaseUrl: cdnBaseUrl}
}

func (s *service) formatPhotoPath(filename string) string {
	if filename == "" {
		return s.cdnBaseUrl + "default.jpg"
	}
	return s.cdnBaseUrl + filename
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

func (s *service) FindOneByAccessToken(c context.Context, accessToken string) (*domain.User, error) {
	user, err := s.repo.FindOneByAccessToken(c, accessToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user.Photo = s.formatPhotoPath(user.Photo)

	return user, nil
}

func (s *service) FindOneByPhone(c context.Context, phone string) (*domain.User, error) {
	user, err := s.repo.FindOneByPhone(c, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	user.Photo = s.formatPhotoPath(user.Photo)

	return user, nil
}

func (s *service) FindOneByID(c context.Context, userId int) (*domain.User, error) {
	user, err := s.repo.FindOneByID(c, userId)
	if err != nil {
		log.Println(err)
		return &domain.User{}, domain.ErrUserNotFound
	}

	user.Photo = s.formatPhotoPath(user.Photo)

	return user, nil
}

func (s *service) FindOneByUsername(c context.Context, username string) (*domain.User, error) {
	user, err := s.repo.FindOneByUsername(c, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.User{}, domain.ErrUserNotFound
		}
		return &domain.User{}, err
	}

	user.Photo = s.formatPhotoPath(user.Photo)

	return user, nil
}

func (s *service) UpdatePassword(c context.Context, userId int, password string) error {
	user, err := s.repo.FindOneByID(c, userId)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	return s.repo.Save(c, user)
}

func (s *service) ChangePassword(c context.Context, userId int, oldPassword, newPassword string) error {
	user, err := s.repo.FindOneByID(c, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	if err := user.ChangePassword(oldPassword, newPassword); err != nil {
		return err
	}

	return s.repo.Save(c, user)
}

func (s *service) UpdatePhoto(c context.Context, userId int, filename string) (*domain.User, error) {
	user, err := s.repo.FindOneByID(c, userId)
	if err != nil {
		return nil, err
	}

	user.SetPhoto(filename)

	if err := s.repo.Save(c, user); err != nil {
		return nil, err
	}

	user.Photo = s.formatPhotoPath(user.Photo)

	return user, nil
}
