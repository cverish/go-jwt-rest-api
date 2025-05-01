package handlers

import (
	"github.com/cverish/go-jwt-rest-api/internal/models"
	"github.com/gin-gonic/gin"
)

// HealthCheck	takes in a gin context and returns an ok response.
//
// Swagger doc autogeneration tags:
//
//	@Summary API Health Check
//	@Tags health check
//	@Success 200 {object} StatusOK "Successful Response"
//	@Router /health [get]
func HealthHandler(c *gin.Context) {
	models.ResponseOK(c, "ok")
}
