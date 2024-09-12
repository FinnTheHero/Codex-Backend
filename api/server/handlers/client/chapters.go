package handlers

import (
	"Codex-Backend/api/models"
	"Codex-Backend/api/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var chapterService = services.NewChapterService()

func FindChapter(c *gin.Context) {
	novel := c.Param("novel")
	title := c.Param("chapter")

	result, err := chapterService.GetChapter(novel, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	chapter, ok := result.(models.Chapter)
	if !ok {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Type assertion failed"},
		)
		return
	}

	c.JSON(http.StatusFound,
		gin.H{"chapter": chapter},
	)
	return
}

func FindAllChapters(c *gin.Context) {
	novel := c.Param("novel")

	result, err := chapterService.GetAllChapters(novel)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	chapters, ok := result.([]models.Chapter)
	if !ok {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Type assertion failed"},
		)
		return
	}

	c.JSON(http.StatusFound,
		gin.H{"chapters": chapters},
	)
	return
}
