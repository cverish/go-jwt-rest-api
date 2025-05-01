package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword converts a plain text password into a hashed version
func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", errors.New("password is too long")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a password against a hash
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword checks password complexity requirements.
// Rules:
//
//	min: 8; max: 64
//	one uppercase character
//	one lowercase character
//	one digit
//	one special character
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(password) > 64 {
		return errors.New("password must be fewer than 64 characters")
	}
	if uppercaseRegex := regexp.MustCompile(`.*?[A-Z]`); !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase character")
	}
	if lowercaseRegex := regexp.MustCompile(`.*?[a-z]`); !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase character")
	}
	if digitRegex := regexp.MustCompile(`.*?[0-9]`); !digitRegex.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	if specialCharRegex := regexp.MustCompile(`.*?[#?!@$%^&*-]`); !specialCharRegex.MatchString(
		password,
	) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
