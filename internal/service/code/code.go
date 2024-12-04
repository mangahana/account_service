package code

import (
	"account/internal/domain"
	"account/internal/infrastructure/configuration"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source ./code.go -destination ./mock/mock.go -package mock
type Repository interface {
	FindOneByPhoneAndIP(c context.Context, phone, ip string) (*domain.Code, error)
	FindLatestByIP(c context.Context, ip string, timestamp time.Time) ([]domain.Code, error)
	FindOneByCredentials(c context.Context, phone, code string) (*domain.Code, error)

	Save(c context.Context, code *domain.Code) error

	RemoveAll(c context.Context, phone string) error
}

type service struct {
	baseUrl     string
	accessToken string

	repo Repository
}

func New(cfg *configuration.SMSConfig, repo Repository) *service {
	return &service{
		baseUrl:     cfg.ApiDomain,
		accessToken: cfg.ApiKey,
		repo:        repo,
	}
}

func (s *service) Send(c context.Context, phone, ip string) error {
	if err := s.spamProtect(c, phone, ip); err != nil {
		return err
	}

	code, err := domain.NewCode(phone, ip)
	if err != nil {
		return err
	}

	if err := s.repo.Save(c, code); err != nil {
		return err
	}

	return s.sendSMS("7"+phone, "mangahana.com\nРастау коды: "+code.Code)
}

func (s *service) spamProtect(c context.Context, phone, ip string) error {
	code, err := s.repo.FindOneByPhoneAndIP(c, phone, ip)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	if code.CreatedAt.UTC().Unix() > time.Now().UTC().Add(time.Minute*-3).Unix() {
		return domain.ErrCodeSendingLimit
	}

	// еще стоит добавить защиту спама на определенный номер,
	// допустим если на один номер за последние пол часа отправилось 5 смс то блокаем смс на этот номер

	lastHalfHour := time.Now().Add(time.Minute * -30)

	codes, err := s.repo.FindLatestByIP(c, ip, lastHalfHour)
	if err != nil {
		return err
	}
	if len(codes) >= 5 {
		return domain.ErrTooManyCodesSent
	}

	return nil
}

func (s *service) sendSMS(recipient, text string) error {
	path := fmt.Sprintf("%s%s?apiKey=%s", s.baseUrl, "/service/Message/SendSmsMessage", s.accessToken)

	data := url.Values{
		"recipient": []string{recipient},
		"text":      []string{text},
	}
	resp, err := http.PostForm(path, data)
	if err != nil {
		return fmt.Errorf("failed to post request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("failed to unmarshal json from body: %v\nbody:%v", err, string(body))
	}

	if res.Code != 0 {
		return errors.New(res.Message)
	}

	return nil
}

func (s *service) Verify(c context.Context, phone, code string) error {
	_, err := s.repo.FindOneByCredentials(c, phone, code)
	return err
}

func (s *service) RemoveAll(c context.Context, phone string) error {
	return s.repo.RemoveAll(c, phone)
}
