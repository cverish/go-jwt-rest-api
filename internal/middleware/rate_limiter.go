package middleware

import (
	"time"

	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware middleware prevents brute force attacks
// by limiting the number of requests to 10 per second.
func RateLimiterMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			models.ResponseTooManyRequests(c, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
