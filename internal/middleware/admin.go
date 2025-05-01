package middleware

import (
	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/cverish/go-jwt-rest-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifies that user is an admin via the database
// if the are no longer an admin, revoke their tokens and force login.
func AdminMiddleware(cfg *config.Config, db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check role from the jwt
		role, exists := c.Get("role")

		if !exists || role != models.RoleAdmin.String() {
			models.ResponseForbidden(c, "admin role required")
			c.Abort()
			return
		}

		// if the jwt claims this is an admin, check their status in the db
		userId, exists := c.Get("user_id")
		if !exists {
			models.ResponseForbidden(c, "admin role required")
			c.Abort()
			return
		}

		user, err := db.GetUserById(userId.(string))
		if err != nil {
			models.ResponseForbidden(c, "admin role required")
			c.Abort()
			return
		}

		if user.Role != models.RoleAdmin {
			// revoke tokens and force login
			utils.RevokeTokens(c, cfg)
			models.ResponseUnauthorized(c, "login required")
			c.Abort()
			return
		}

		c.Next()
	}
}
