package helper

import (
	"Codex-Backend/api/internal/domain"

	"github.com/gin-gonic/gin"
)

// Helper function to get user from context
func GetUserFromContext(c *gin.Context) (*domain.User, bool) {
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*domain.User); ok {
			return u, true
		}
	}
	return nil, false
}
