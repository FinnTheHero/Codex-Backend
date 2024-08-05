package client_handler

import (
	aws_methods "Codex-Backend/api/aws/methods"

	"github.com/gin-gonic/gin"
)

func FindNovel(c *gin.Context) {
	title := c.Param("novel")

	NovelSchema, err := aws_methods.GetNovel(title)
	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"novel": NovelSchema.Novel,
	})
}

func FindAllNovels(c *gin.Context) {
	Novels, err := aws_methods.GetAllNovels()
	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"novels": Novels,
	})
}
