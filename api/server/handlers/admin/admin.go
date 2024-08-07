package admin_handler

import (
	aws_methods "Codex-Backend/api/aws/methods"
	"Codex-Backend/api/types"
	"Codex-Backend/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNovel(c *gin.Context) {
	var novel types.Novel

	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := aws_methods.CreateTable(novel.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = aws_methods.CreateNovel(novel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Novel created successfully",
	})
}

func CreateChapter(c *gin.Context) {
	novelTitle := c.Param("novel")

	var chapter types.Chapter

	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tableExists, err := utils.IsTableCreated(novelTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !tableExists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Novel not found",
		})
		return
	}

	err = aws_methods.CreateChapter(novelTitle, chapter)
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
