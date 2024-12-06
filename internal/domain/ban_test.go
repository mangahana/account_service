package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ban, err := NewBan(1, 2, "spam", time.Now().Add(time.Second*5))

		assert.NoError(t, err)
		assert.NotZero(t, ban)

		t.Log(ban)
	})
}
