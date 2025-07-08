package middleware

import (
	"Codex-Backend/api/internal/infrastructure/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyUsersTablesExist() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := repository.VerifyUsersTable()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Next()
	}
}

func VerifyNovelsTableExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := repository.VerifyNovelsTable()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Next()
	}
}
