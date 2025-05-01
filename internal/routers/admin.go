package routers

import (
	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/handlers"
	"github.com/cverish/go-jwt-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AttachAdminRouterGroup attaches the middleware and routes associated with admin authorization
// to the gin router.
func AttachAdminRouterGroup(cfg *config.Config, db *database.Database, router *gin.RouterGroup) {
	adminHandler := handlers.NewAdminHandler(db)

	public := router.Group("/admin")
	{
		public.POST("/manage/create-initial-admin", adminHandler.CreateInitialAdmin)
	}

	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg, db))
	admin.Use(middleware.AdminMiddleware(cfg, db))
	{
		admin.GET("/invites", adminHandler.GetInvitedUsers)
		admin.POST("/invites/create", adminHandler.CreateInvite)
		admin.POST("/invites/reset-password", adminHandler.ResetInvitedUserPassword)
		admin.DELETE("/invites/delete/:user_id", adminHandler.DeleteInvitedUser)
		admin.GET("/users", adminHandler.GetUsers)
		admin.POST("/users/reset-password", adminHandler.ResetUserPassword)
		admin.DELETE("/users/delete/:user_id", adminHandler.DeleteUser)
		admin.POST("/manage/add-admin/:user_id", adminHandler.AddAdmin)
		admin.POST("/manage/remove-admin/:user_id", adminHandler.RemoveAdmin)
	}
}
