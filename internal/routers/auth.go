package routers

import (
	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/handlers"
	"github.com/cverish/go-jwt-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AttachAuthRouterGroup attaches the routes and middleware associated with authentication.
func AttachAuthRouterGroup(cfg *config.Config, db *database.Database, router *gin.RouterGroup) {
	authHandler := handlers.NewAuthHandler(db, cfg)

	// public routes
	public := router.Group("/auth")
	{
		public.POST("/login", authHandler.Login)
	}
	// protected routes
	protected := router.Group("/auth")
	protected.Use(middleware.AuthMiddleware(cfg, db))
	{
		protected.POST("/refresh-token", authHandler.RefreshToken)
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.POST("/logout", authHandler.Logout)
	}
}
