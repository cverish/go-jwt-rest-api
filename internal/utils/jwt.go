package utils

import (
	"errors"
	"time"

	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateTokenString(jwtSecret []byte, jwtExpiration time.Duration, now time.Time, userId string, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetTokenClaims(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		} else {
			return nil, errors.New("invalid or expired token")
		}
	}
	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}
	return claims, nil
}

// sets the given tokenName in the gin context with the given accessToken and expiry values
func SetToken(c *gin.Context, cfg *config.Config, tokenName string, accessToken string, expiry time.Duration) {
	c.SetCookie(
		tokenName,
		accessToken,
		int(expiry),
		"/",
		cfg.JWT.CookieDomain,
		!cfg.IsDev,
		true,
	)
}

// revoke access and refresh tokens
func RevokeTokens(c *gin.Context, cfg *config.Config) {
	c.SetCookie(cfg.JWT.AccessTokenKey, "", -1, "/", cfg.JWT.CookieDomain, !cfg.IsDev, true)
	c.SetCookie(cfg.JWT.RefreshTokenKey, "", -1, "/", cfg.JWT.CookieDomain, !cfg.IsDev, true)
}
