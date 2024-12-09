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
	ErrInvalidPassword      = errors.New("INVALID_PASSWORD")
	ErrUnsupportedFormat    = errors.New("UNSUPPORTED_FORMAT")

	// bans
	ErrBanNotFound         = errors.New("BAN_NOT_FOUND")
	ErrCantBanYourself     = errors.New("CANT_BAN_YOURSELF")
	ErrReasonCantBeEmpty   = errors.New("REASON_CANT_BE_EMPTY")
	ErrExpiryCantBePast    = errors.New("EXPIRY_CANT_BE_PAST")
	ErrCannotBanHigherUser = errors.New("CANT_BAN_HIGHER_USER")
)
