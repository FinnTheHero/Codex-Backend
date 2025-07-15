package handlers

import (
	"Codex-Backend/api/internal/domain"
	chapter_service "Codex-Backend/api/internal/usecases/chapter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindChapter(c *gin.Context) {
	novelId := c.Param("novel")
	titleId := c.Param("chapter")

	chapter, err := chapter_service.GetChapter(novelId, titleId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chapter": chapter,
	})
}

func FindPreviousAndNextChapters(c *gin.Context) {
	novelId := c.Param("novel")
	chapterId := c.Param("chapter")

	chapters, err := chapter_service.GetAllChapters(novelId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	prev_next := []domain.Chapter{}

	for i, chapter := range chapters {
		if chapter.Title == chapterId {
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

	c.JSON(http.StatusOK, gin.H{
		"chapters": prev_next,
	})
}

func FindAllChapters(c *gin.Context) {
	novelId := c.Param("novel")

	chapters, err := chapter_service.GetAllChapters(novelId)
	if err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chapters": chapters,
	})
}
