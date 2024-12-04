package domain

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int
	Username    string
	Phone       string
	Password    string
	Photo       *string
	Description *string
	Role        *Role
	CreatedAt   time.Time
}

func NewUser(username, phone, password string) (*User, error) {
	row := regexp.MustCompile("^[a-zA-Z0-9]+(_?[a-zA-Z0-9]+)*$")
	if !row.Match([]byte(username)) {
		return nil, ErrInvalidUsername
	}

	if len(username) < 3 || len(username) > 25 {
		return nil, ErrInvalidUsername
	}

	if len(password) < 8 {
		return nil, ErrTooShortPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: username,
		Phone:    phone,
		Password: string(hashedPassword),
	}, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
