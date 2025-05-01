package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("Password hashes correctly", func(t *testing.T) {
		hashed_pw, err := HashPassword("my-password")
		require.NotEqual(t, hashed_pw, "", "hashed password should not be empty string")
		require.Nil(t, err, "should not have an error")
	})

	t.Run("Password does not hash if too long", func(t *testing.T) {
		long_pw := "a_password_with_more_than_72_characters_it_is_very_very_long_as_you_can_see"
		hashed_pw, err := HashPassword(long_pw)

		require.Equal(t, hashed_pw, "", "hashed password should be empty")
		require.Errorf(t, err, "password is too long", "should have error")
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("Password hashes match", func(t *testing.T) {
		hashed_pw, _ := bcrypt.GenerateFromPassword([]byte("my-password"), 12)
		res := CheckPasswordHash("my-password", string(hashed_pw))
		require.True(t, res, "password hashes should match")
	})

	t.Run("Password hashes do not match", func(t *testing.T) {
		hashed_pw, _ := bcrypt.GenerateFromPassword([]byte("my-password"), 12)
		res := CheckPasswordHash("my-wrong-password", string(hashed_pw))
		require.False(t, res, "password hashes should not match")
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("Password is valid", func(t *testing.T) {
		pw := "myValidpassword123!"
		require.Nil(t, ValidatePassword(pw), "password should be valid")
	})

	t.Run("Password is too short", func(t *testing.T) {
		pw := "Short!1"
		require.Errorf(t, ValidatePassword(pw), "password must be at least 8 characters", "should return err")
	})

	t.Run("Password is too long", func(t *testing.T) {
		pw := "A_password_with_more_than_72_characters_it_is_very_very_long_as_you_can_see!"
		require.Errorf(t, ValidatePassword(pw), "password must be fewer than 64 characters", "should return err")
	})

	t.Run("Password does not contain uppercase character", func(t *testing.T) {
		pw := "a_lowercase123!"
		require.Errorf(t, ValidatePassword(pw), "password must contain at least one uppercase character", "should return err")
	})

	t.Run("Password does not contain lowercase character", func(t *testing.T) {
		pw := "AN_UPPERCASE123!"
		require.Errorf(t, ValidatePassword(pw), "password must contain at least one lowercase character", "should return err")
	})

	t.Run("Password does not contain digit", func(t *testing.T) {
		pw := "No_digits!"
		require.Errorf(t, ValidatePassword(pw), "password must contain at least one digit", "should return err")
	})

	t.Run("Password does not contain special character", func(t *testing.T) {
		pw := "Nospecialchars1"
		require.Errorf(t, ValidatePassword(pw), "password must contain at least one special character", "should return err")
	})
}
