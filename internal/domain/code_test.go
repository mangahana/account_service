package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		code, err := NewCode("7778881133", "127.0.0.1")
		assert.NoError(t, err)
		assert.NotZero(t, code)
	})
}
