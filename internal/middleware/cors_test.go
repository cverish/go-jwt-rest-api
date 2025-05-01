package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCorsMiddleware(t *testing.T) {
	t.Run("CORS headers set correctly by middleware", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "http://example.com", nil)

		CorsMiddleware()(c)

		require.Equal(t, c.Writer.Header().Get("Access-Control-Allow-Origin"), "localhost")
		require.Equal(
			t,
			c.Writer.Header().Get("Access-Control-Allow-Methods"),
			"POST, GET, OPTIONS, PUT, DELETE",
		)
		require.Equal(
			t,
			c.Writer.Header().Get("Access-Control-Allow-Headers"),
			"Content-Type, Authorization, Cookie",
		)
	})

	t.Run("OPTIONS method aborts", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("OPTIONS", "http://example.com", nil)

		CorsMiddleware()(c)

		require.True(t, c.IsAborted())
		require.Equal(t, c.Writer.Status(), http.StatusNoContent)
	})
}
