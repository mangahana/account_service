package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	type testCase struct {
		name     string
		username string
		phone    string
		password string
		wantErr  error
	}

	testCases := []testCase{
		{
			name:     "success",
			username: "john",
			phone:    "8887779900",
			password: "12345678",
			wantErr:  nil,
		},
		{
			name:     "incorrect username",
			username: "john_",
			phone:    "8887779900",
			password: "12345678",
			wantErr:  ErrInvalidUsername,
		},
		{
			name:     "incorrect password",
			username: "john",
			phone:    "8887779900",
			password: "short",
			wantErr:  ErrTooShortPassword,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.username, tt.phone, tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}
