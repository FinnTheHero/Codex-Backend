package handlers

import (
	"Codex-Backend/api/internal/domain"
	chapter_service "Codex-Backend/api/internal/usecases/chapter"
	novel_service "Codex-Backend/api/internal/usecases/novel"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNovel(c *gin.Context) {
	var novel domain.Novel

	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := novel_service.CreateNovel(novel)
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

	err := chapter_service.CreateChapter(novel, chapter)
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
