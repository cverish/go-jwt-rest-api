package middleware

import (
	"errors"
	"time"

	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT tokens in incoming requests.
func AuthMiddleware(cfg *config.Config, db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(cfg.JWT.AccessTokenKey)

		if err != nil {
			models.ResponseUnauthorized(c, "missing cookie")
			c.Abort()
			return
		}

		claims, err := utils.GetTokenClaims(cookie, []byte(cfg.JWT.AccessTokenSecret))
		// access token is still valid
		if err == nil {
			// Set user information in context
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
			c.Next()
			return
		} else if !cfg.JWT.BackendRefresh {
			models.ResponseUnauthorized(c, "login required")
			c.Abort()
			return
		}

		// try refreshing access token from refresh token
		if err := refreshAccessToken(cfg, db, c); err != nil {
			models.ResponseUnauthorized(c, "login required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// refreshAccessToken checks the validity of the refresh token and
// issues a new access token if refresh token is valid.
//
// Returns:
//
//	nil: success
//	error: invalid token or other error
func refreshAccessToken(cfg *config.Config, db *database.Database, c *gin.Context) error {
	// check refresh token
	refreshCookie, err := c.Cookie(cfg.JWT.RefreshTokenKey)
	if err != nil {
		return err
	}

	// extract claims
	claims, err := utils.GetTokenClaims(refreshCookie, []byte(cfg.JWT.RefreshTokenSecret))
	if err != nil {
		return err
	}

	// validate that user is still authorized
	user, _ := db.GetUserById(claims["user_id"].(string))
	if user == nil {
		return errors.New("user is no longer authorized")
	}

	// generate new access token
	accessToken, err := utils.GenerateTokenString(
		[]byte(cfg.JWT.AccessTokenSecret),
		cfg.JWT.AccessTokenExpiry,
		time.Now(),
		user.ID.String(),
		user.Email,
		user.Role.String(),
	)
	if err != nil {
		return err
	}

	utils.SetToken(c, cfg, cfg.JWT.AccessTokenKey, accessToken, cfg.JWT.AccessTokenExpiry)

	c.Set("user_id", user.ID.String())
	c.Set("email", user.Email)
	c.Set("role", user.Role.String())
	return nil
}
