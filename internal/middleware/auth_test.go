package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"github.com/cverish/go-jwt-rest-api/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("Correct token passes through middleware", func(t *testing.T) {
		db, _, err := testutils.GetMockGormDB()
		if err != nil {
			t.Fatal("error initializing db")
		}
		cfg, err := testutils.GetMockConfig()
		if err != nil {
			t.Fatal("error getting config")
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		now := time.Now()
		jwtSecret := cfg.JWT.AccessTokenSecret
		tokenString, _ := utils.GenerateTokenString(
			[]byte(jwtSecret),
			cfg.JWT.AccessTokenExpiry,
			now,
			"12345",
			"email@example.com",
			"user",
		)

		cookie := http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			Path:     "/",
			MaxAge:   int(time.Duration(10 * time.Minute)),
			HttpOnly: true,
			Secure:   true,
		}
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)
		c.Request.AddCookie(&cookie)

		AuthMiddleware(cfg, &database.Database{DB: db})(c)

		require.False(t, c.IsAborted())
		require.Equal(t, c.MustGet("user_id"), "12345")
		require.Equal(t, c.MustGet("email"), "email@example.com")
		require.Equal(t, c.MustGet("role"), "user")
	})

	t.Run("Missing token does not pass through middleware", func(t *testing.T) {
		db, _, err := testutils.GetMockGormDB()
		if err != nil {
			t.Fatal("error initializing db")
		}
		cfg, err := testutils.GetMockConfig()
		if err != nil {
			t.Fatal("error getting config")
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)
		c.Request.Header.Add("Authorization", "abc 1")

		AuthMiddleware(cfg, &database.Database{DB: db})(c)
		require.True(t, c.IsAborted())
		require.Equal(t, c.Writer.Status(), 401)
	})

	t.Run("Incorrect token does not pass through middleware", func(t *testing.T) {
		db, _, err := testutils.GetMockGormDB()
		if err != nil {
			t.Fatal("error initializing db")
		}
		cfg, err := testutils.GetMockConfig()
		if err != nil {
			t.Fatal("error getting config")
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)
		c.Request.Header.Add("Authorization", "abc 1")

		now := time.Now()
		tokenString, _ := utils.GenerateTokenString(
			[]byte("different_secret"),
			time.Duration(10*time.Minute),
			now,
			"12345",
			"email@example.com",
			"user",
		)

		cookie := http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			Path:     "/",
			MaxAge:   int(time.Duration(10 * time.Minute)),
			HttpOnly: true,
			Secure:   true,
		}
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)
		c.Request.AddCookie(&cookie)

		AuthMiddleware(cfg, &database.Database{DB: db})(c)
		require.True(t, c.IsAborted())
		require.Equal(t, c.Writer.Status(), 401)
	})
}
