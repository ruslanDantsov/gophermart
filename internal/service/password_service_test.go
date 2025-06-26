package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordService_Compare(t *testing.T) {
	ps := &PasswordService{}

	t.Run("successfully compares matching password", func(t *testing.T) {
		password := "testPassword"
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		err = ps.Compare(string(hashed), password)
		assert.NoError(t, err)
	})

	t.Run("fails on wrong password", func(t *testing.T) {
		password := "correctPassword"
		wrongPassword := "wrongPassword"
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		err = ps.Compare(string(hashed), wrongPassword)
		assert.Error(t, err)
		assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
	})

	t.Run("fails on invalid hash", func(t *testing.T) {
		invalidHash := "not-a-valid-bcrypt-hash"
		err := ps.Compare(invalidHash, "anyPassword")
		assert.Error(t, err)
	})
}
