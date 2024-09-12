package admin_handler

import (
	"Codex-Backend/api/models"
	admin_services "Codex-Backend/api/server/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

var adminService = admin_services.NewAdminService()

func CreateNovel(c *gin.Context) {
	var novel models.Novel

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

	var chapter models.Chapter

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
