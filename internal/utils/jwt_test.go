package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestGenerateTokenString(t *testing.T) {
	jwtSecret := []byte("my-secret")
	jwtExpiration := time.Duration(5 * time.Minute)
	now := time.Now()
	userId := "my-user-id"
	email := "my-email@example.com"
	role := "user"

	expected_claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(jwtExpiration).Unix(),
	}
	expected_token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, expected_claims).
		SignedString(jwtSecret)

	t.Run("Token generates without error", func(t *testing.T) {
		res, err := GenerateTokenString(jwtSecret, jwtExpiration, now, userId, email, role)
		require.NotNil(t, res, "result should not be nil")
		require.Nil(t, err, "err should be nil")
	})

	t.Run("Expected token generated", func(t *testing.T) {
		res, _ := GenerateTokenString(jwtSecret, jwtExpiration, now, userId, email, role)
		require.Equal(t, expected_token, res, "generated tokens should be equal")
	})

	t.Run("Expected claims added to token", func(t *testing.T) {
		res, _ := GenerateTokenString(jwtSecret, jwtExpiration, now, userId, email, role)

		token, _ := jwt.Parse(
			res,
			func(token *jwt.Token) (interface{}, error) { return jwtSecret, nil },
		)
		claims, _ := token.Claims.(jwt.MapClaims)

		require.Equal(t, claims["user_id"], userId, "user ID should be equal")
		require.Equal(t, claims["email"], email, "email should be equal")
		require.Equal(t, claims["role"], role, "roles should be equal")
	})
}

func TestGetTokenClaims(t *testing.T) {
	jwtSecret := []byte("my-secret")
	jwtExpiration := time.Duration(5 * time.Second)
	now := time.Now()
	userId := "my-user-id"
	email := "my-email@example.com"
	role := "user"

	expected_claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(jwtExpiration).Unix(),
	}
	expected_token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, expected_claims).
		SignedString(jwtSecret)

	t.Run("token claims should be returned", func(t *testing.T) {
		token_claims, err := GetTokenClaims(expected_token, jwtSecret)

		require.NotNil(t, token_claims, "token claims should not be nil")
		require.Nil(t, err, "error should be nil")
	})

	t.Run("error should be returned if token is empty string", func(t *testing.T) {
		token_claims, err := GetTokenClaims("", jwtSecret)

		require.Nil(t, token_claims, "token claims should be nil")
		require.Errorf(t, err, "authorization header missing", "error should be returned")
	})

	t.Run("error should be returned if incorrect signing method", func(t *testing.T) {
		token, _ := jwt.NewWithClaims(jwt.SigningMethodES256, expected_claims).
			SignedString(jwtSecret)
		token_claims, err := GetTokenClaims(token, jwtSecret)

		require.Nil(t, token_claims, "token claims should be nil")
		require.Errorf(t, err, "signature is invalid", "error should be returned")
	})

	t.Run("error should be returned if token expired", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": userId,
			"email":   email,
			"role":    role,
			"iat":     now.Unix(),
			"exp":     now.Add(-jwtExpiration).Unix(),
		}
		token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
		token_claims, err := GetTokenClaims(token, jwtSecret)

		require.Nil(t, token_claims, "token claims should be nil")
		require.Errorf(t, err, "token expired", "error should be returned")
	})

}
