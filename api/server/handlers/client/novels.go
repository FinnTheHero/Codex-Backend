package handlers

import (
	"Codex-Backend/api/models"
	client_services "Codex-Backend/api/server/services/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

var novelService = client_services.NewNovelService()

func FindNovel(c *gin.Context) {
	title := c.Param("novel")

	result, err := novelService.GetNovel(title)
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

	novel, ok := result.(models.NovelDTO)
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
	result, err := novelService.GetAllNovels()
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

	novels, ok := result.([]models.NovelDTO)
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
