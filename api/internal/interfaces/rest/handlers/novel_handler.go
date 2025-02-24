package handlers

import (
	"Codex-Backend/api/internal/domain"
	novel_service "Codex-Backend/api/internal/usecases/novel"

	"net/http"

	"github.com/gin-gonic/gin"
)



func FindNovel(c *gin.Context) {
	title := c.Param("novel")

	result, err := novel_service.GetNovel(title)
	if err != nil {
		errStatus := http.StatusInternalServerError

		if err.Error() == (title + " not found") {
			errStatus = http.StatusNotFound
		}

		c.JSON(errStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	novel, ok := result.(domain.NovelDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Type assertion failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novel": novel.Novel,
	})
	return
}

func FindAllNovels(c *gin.Context) {
	result, err := novel_service.GetAllNovels()
	if err != nil {
		errStatus := http.StatusInternalServerError

		if err.Error() == "No novels found" {
			errStatus = http.StatusNotFound
		}

		c.JSON(errStatus, gin.H{
			"error": err.Error(),
		})
		return
	}

	novels, ok := result.([]domain.Novel)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Type assertion failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novels": novels,
	})
	return
}
