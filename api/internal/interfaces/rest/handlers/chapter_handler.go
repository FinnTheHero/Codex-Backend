package handlers

import (
	"Codex-Backend/api/internal/domain"
	chapter_service "Codex-Backend/api/internal/usecases/chapter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindChapter(c *gin.Context) {
	novel := c.Param("novel")
	title := c.Param("chapter")

	chapter, err := chapter_service.GetChapter(novel, title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		},
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

	chapters, err := chapter_service.GetAllChapters(novel)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	prev_next := []domain.Chapter{}

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

	chapters, err := chapter_service.GetAllChapters(novel)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"chapters": chapters},
	)
	return
}
