package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	// TODO: Add better health check for future update.
	// Implement resource and status monitoring.
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
