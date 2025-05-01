package routers

import (
	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/handlers"
	"github.com/cverish/go-jwt-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AttachUserRouterGroup attaches the routes and middleware associated with users.
func AttachUserRouterGroup(cfg *config.Config, db *database.Database, router *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler(db)

	// public routes
	public := router.Group("/")
	public.POST("/register", userHandler.Register)

	// protected routes
	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware(cfg, db))
	{
		protected.GET("/:user_id", userHandler.GetUser)
		protected.PUT("/:user_id", userHandler.UpdateUser)
	}
}
