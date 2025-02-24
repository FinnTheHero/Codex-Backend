package handlers

import (
	"Codex-Backend/api/internal/domain"
	user_service "Codex-Backend/api/internal/usecases/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

var adminService = user_service.NewAdminService()

func CreateNovel(c *gin.Context) {
	var novel domain.Novel

	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := adminService.CreateNovel(novel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Novel created successfully",
	})
}

func CreateChapter(c *gin.Context) {
	novel := c.Param("novel")

	var chapter domain.Chapter

	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := adminService.CreateChapter(novel, chapter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Chapter created successfully",
	})
}
