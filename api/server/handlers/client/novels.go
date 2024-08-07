package client_handler

import (
	aws_methods "Codex-Backend/api/aws/methods"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindNovel(c *gin.Context) {
	title := c.Param("novel")

	NovelSchema, err := aws_methods.GetNovel(title)
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

	c.JSON(http.StatusFound, gin.H{
		"novel": NovelSchema.Novel,
	})
}

func FindAllNovels(c *gin.Context) {
	Novels, err := aws_methods.GetAllNovels()
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

	c.JSON(http.StatusFound, gin.H{
		"novels": Novels,
	})
}
