package server

import (
	"log"
	"net/http"

	_ "github.com/cverish/go-jwt-rest-api/docs"
	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/handlers"
	"github.com/cverish/go-jwt-rest-api/internal/middleware"
	"github.com/cverish/go-jwt-rest-api/internal/routers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHttpServer takes in the app config and database, sets up app-level middleware,
// and returns the API http server.
func NewHttpServer(cfg *config.Config, db *database.Database) *http.Server {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	apiGroup := r.Group(cfg.API.BasePath)

	// attach docs before middleware if in dev
	if cfg.IsDev {
		apiGroup.GET("/docs/*any", func(c *gin.Context) {
			// redirect to index page
			if c.Param("any") == "" || c.Param("any") == "/" {
				c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
				return
			}
			ginSwagger.WrapHandler(
				swaggerFiles.Handler,
				ginSwagger.URL(cfg.API.BasePath+"/docs/doc.json"),
				ginSwagger.DefaultModelsExpandDepth(2),
			)(c)
		})
	}

	// attach CORS and Security middleware to all routes
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.SecurityMiddleware(cfg.IsDev))

	// attach health route
	apiGroup.GET("/health", handlers.HealthHandler)

	// attach router groups
	routers.AttachAuthRouterGroup(cfg, db, apiGroup)
	routers.AttachUserRouterGroup(cfg, db, apiGroup)
	routers.AttachAdminRouterGroup(cfg, db, apiGroup)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)

	return &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}
