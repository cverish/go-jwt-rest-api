package middleware

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSecurity(t *testing.T) {
	t.Run("Non-local dev should not allow http", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)

		SecurityMiddleware(false)(c)

		require.True(t, c.IsAborted())
	})

	t.Run("Headers are added to request", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "https://example.com", nil)

		SecurityMiddleware(false)(c)

		// we'll check the nonce in another test
		require.NotEqual(t, c.Writer.Header().Get("Content-Security-Policy"), "")
		require.Equal(t, c.Writer.Header().Get("X-Frame-Options"), "DENY")
		require.Equal(t, c.Writer.Header().Get("X-Content-Type-Options"), "nosniff")
	})

	t.Run("Nonce is in header and matches context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "https://example.com", nil)

		SecurityMiddleware(false)(c)

		nonce, exists := c.Get("csp_nonce")
		require.True(t, exists, "CSP nonce should be written to context")
		require.Equal(
			t,
			c.Writer.Header().Get("Content-Security-Policy"),
			fmt.Sprintf("script-src 'nonce-%s'", nonce),
			"Content-Security-Policy nonce should match nonce in context",
		)
	})

	t.Run("Redirects should not have headers added", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "https://example.com", nil)
		c.Status(302)

		SecurityMiddleware(false)(c)

		require.Equal(t, c.Writer.Header().Get("Content-Security-Policy"), "")
		require.Equal(t, c.Writer.Header().Get("X-Frame-Options"), "")
		require.Equal(t, c.Writer.Header().Get("X-Content-Type-Options"), "")
	})

	t.Run("Local dev should allow http", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)

		SecurityMiddleware(true)(c)

		require.False(t, c.IsAborted())
	})
}
