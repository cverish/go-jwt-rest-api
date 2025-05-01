package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// adds 200 StatusOK response with message to gin context
func ResponseOK(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// adds 200 StatusOK response with item to gin context
func ResponseOKItem(c *gin.Context, item any) {
	c.JSON(http.StatusOK, gin.H{"item": item})
}

// adds 200 StatusOK response with items list to gin context
func ResponseOKList(c *gin.Context, items any) {
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// adds 201 StatusCreated response with message and id of created item to gin context
func ResponseCreated(c *gin.Context, message string, id string) {
	c.JSON(http.StatusCreated, gin.H{"message": message, "id": id})
}

// adds 400 StatusBadRequest response with error message to gin context
func ResponseBadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// adds 401 StatusUnauthorized response with error message to gin context
func ResponseUnauthorized(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": err})
}

// adds 403 StatusForbidden response with error message to gin context
func ResponseForbidden(c *gin.Context, err string) {
	c.JSON(http.StatusForbidden, gin.H{"error": err})
}

// adds 404 StatusNotFound response with error message to gin context
func ResponseNotFound(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, gin.H{"error": err})
}

// adds 406 StatusNotAcceptable response with error message to gin context
func ResponseNotAcceptable(c *gin.Context, err string) {
	c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
}

// adds 429 TooManyRequests response with error message to gin context
func ResponseTooManyRequests(c *gin.Context, err string) {
	c.JSON(http.StatusTooManyRequests, gin.H{"error": err})
}

// adds 500 StatusInternalServerError response with generic error message to gin context
func ResponseInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
