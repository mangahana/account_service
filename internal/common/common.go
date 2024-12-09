package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateRandomHash() (string, error) {
	randomData := make([]byte, 256)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(randomData)
	return hex.EncodeToString(hash[:]), nil
}
