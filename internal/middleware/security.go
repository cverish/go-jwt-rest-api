package middleware

import (
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func SecurityMiddleware(isDev bool) gin.HandlerFunc {
	opt := secure.Options{
		ContentSecurityPolicy: "default-src 'self'; script-src $NONCE; object-src 'none';",
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		SSLRedirect:           true,
		STSSeconds:            31536000,
		IsDevelopment:         isDev,
	}

	return func(c *gin.Context) {
		// avoid header rewrite if response is a redirection and continue to next handler
		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Next()
			return
		}

		secureMiddleware := secure.New(opt)

		nonce, err := secureMiddleware.ProcessAndReturnNonce(c.Writer, c.Request)
		if err != nil {
			models.ResponseInternalServerError(c)
			c.Abort()
			return
		}

		c.Set("csp_nonce", nonce)
		c.Next()
	}
}
