package domain

import (
	"crypto/rand"
	"math/big"
	"time"
)

type Code struct {
	Code      string
	Phone     string
	IP        string
	CreatedAt time.Time
}

func randomCode(length int) (string, error) {
	var output string
	str := "1234567890"

	for i := 0; i < length; i++ {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(len(str)-1)))
		if err != nil {
			return "", err
		}
		output = output + string(str[i.Int64()])
	}

	return output, nil
}

func NewCode(phone, ip string) (*Code, error) {
	code, err := randomCode(4)
	if err != nil {
		return nil, err
	}

	return &Code{
		Code:      code,
		Phone:     phone,
		IP:        ip,
		CreatedAt: time.Now().UTC(),
	}, nil
}
