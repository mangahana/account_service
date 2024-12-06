package domain

import "time"

type Ban struct {
	ID         int
	BannedByID int
	UserID     int
	Reason     string
	Expiry     time.Time
	IsActive   bool
	CreatedAt  time.Time
}

func NewBan(bannedBy, userID int, reason string, expiry time.Time) (*Ban, error) {
	if bannedBy == userID {
		return &Ban{}, ErrCantBanYourself
	}

	if len(reason) == 0 {
		return &Ban{}, ErrReasonCantBeEmpty
	}

	if time.Now().After(expiry) {
		return &Ban{}, ErrExpiryCantBePast
	}

	return &Ban{
		BannedByID: bannedBy,
		UserID:     userID,
		Reason:     reason,
		Expiry:     expiry,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}, nil
}
