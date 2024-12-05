package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("USER_NOT_FOUND")
	ErrPhoneNotFound        = errors.New("PHONE_NOT_FOUND")
	ErrInvalidUsername      = errors.New("INVALID_USERNAME")
	ErrTooShortPassword     = errors.New("TOO_SHORT_PASSWORD")
	ErrPhoneAlreadyInUse    = errors.New("PHONE_ALREADY_IN_USE")
	ErrTooManyCodesSent     = errors.New("TOO_MANY_CODES_SENT")
	ErrUsernameAlreadyInUse = errors.New("USERNAME_ALREADY_IN_USE")
	ErrCodeSendingLimit     = errors.New("CODE_SENDING_LIMIT")
	ErrInvalidCredentials   = errors.New("INVALID_CREDENTIALS")
)
