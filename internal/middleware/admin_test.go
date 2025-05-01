package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAdminMiddleware(t *testing.T) {
	t.Run("Admin users pass through middleware", func(t *testing.T) {
		db, _, err := testutils.GetMockGormDB()
		if err != nil {
			t.Fatal("error initializing db")
		}
		cfg, err := testutils.GetMockConfig()
		if err != nil {
			t.Fatal("error getting config")
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("role", "admin")

		AdminMiddleware(cfg, &database.Database{DB: db})(c)
		require.False(t, c.IsAborted())
	})

	t.Run("Regular users do not pass through middleware", func(t *testing.T) {
		db, _, err := testutils.GetMockGormDB()
		if err != nil {
			t.Fatal("error initializing db")
		}
		cfg, err := testutils.GetMockConfig()
		if err != nil {
			t.Fatal("error getting config")
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("role", "user")

		AdminMiddleware(cfg, &database.Database{DB: db})(c)
		require.True(t, c.IsAborted())
	})

	t.Run("Anonymous role does not pass through middleware", func(t *testing.T) {
		db, _, err := testutils.GetMockGormDB()
		if err != nil {
			t.Fatal("error initializing db")
		}
		cfg, err := testutils.GetMockConfig()
		if err != nil {
			t.Fatal("error getting config")
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		AdminMiddleware(cfg, &database.Database{DB: db})(c)
		require.True(t, c.IsAborted())
	})
}
