package dtos

import (
	"account/internal/domain"
	"time"
)

type RegisterInput struct {
	Phone string `json:"phone"`
	IP    string // client ip address
}

type ConfirmCodeInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type CompleteRegisterInput struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInput struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AuthOutput struct {
	AccessToken string `json:"access_token"`
}

type RecoveryInput struct {
	Phone string `json:"phone"`
	IP    string
}

type CompleteRecovery struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type BanInput struct {
	CallerUserID int
	UserID       int
	Reason       string
	Expiry       time.Time
}

type UnBanInput struct {
	UserID int
	BanID  int
	Reason string
}

type AuthenticateInput struct {
	AccessToken string
}

type AuthenticateOutput struct {
	UserID   int
	IsBanned bool
	Role     domain.Role
}

type UserOutput struct {
	ID          int
	Username    string
	Photo       string
	Description string
	Role        domain.Role
}

func NewUserOutput(id int, username, photo, description string, role domain.Role) *UserOutput {
	return &UserOutput{
		ID:          id,
		Username:    username,
		Photo:       photo,
		Description: description,
		Role:        role,
	}
}

type FindByUsernameInput struct {
	Username string
}

type FindByIDInput struct {
	ID int
}

type IsPhoneExistsInput struct {
	Phone string
}

type ChangePasswordInput struct {
	UserID      int
	OldPassword string
	NewPassword string
	Logout      bool
}

type UploadImageInput struct {
	UserID int
	File   []byte
}

type UploadImageOutput struct {
	Filename string
}

type RemovePhotoInput struct {
	UserID int
}
