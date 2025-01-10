package handlers

import (
	"Codex-Backend/api/models"
	client_services "Codex-Backend/api/server/services/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

var chapterService = client_services.NewChapterService()

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

	c.JSON(http.StatusOK,
		gin.H{"chapter": chapter},
	)
	return
}

func FindPreviousAndNextChapters(c *gin.Context) {
	novel := c.Param("novel")
	title := c.Param("chapter")

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

	prev_next := []models.Chapter{}

	for i, chapter := range chapters {
		if chapter.Title == title {
			if i > 0 {
				prev_next = append(prev_next, chapters[i-1])
			}
			prev_next = append(prev_next, chapter)
			if i < len(chapters)-1 {
				prev_next = append(prev_next, chapters[i+1])
			}
			break
		}
	}

	c.JSON(http.StatusOK,
		gin.H{"chapters": prev_next},
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

	c.JSON(http.StatusOK,
		gin.H{"chapters": chapters},
	)
	return
}
