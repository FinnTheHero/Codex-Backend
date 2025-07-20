package firestore_handlers

import (
	"Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_services "Codex-Backend/api/internal/usecases-firestore/collections"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindChapter(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")
	chapterId := c.Param("chapter")

	if novelId == "" || chapterId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "IDs are not present in request",
		})
		return
	}

	chapter, err := firestore_services.GetChapter(novelId, chapterId, ctx)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to retrieve chapter: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve chapter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chapter": chapter,
	})
}

func FindAllChapters(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")

	if novelId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel ID not found",
		})
		return
	}

	chapters, err := firestore_services.GetAllChapters(novelId, ctx)
	if e, ok := err.(*common.Error); ok {
		if e.StatusCode() == 404 {
			c.AbortWithStatusJSON(e.StatusCode(), gin.H{
				"error": "Failed to retrieve chapters: " + e.Error(),
			})
			return
		}
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve chapters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chapters": chapters,
	})
}

func CreateChapter(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")

	chapter := domain.Chapter{}

	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get chapter data: " + err.Error(),
		})
		return
	}

	err := firestore_services.CreateChapter(novelId, chapter, ctx)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to create chapter: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create chapter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"chapter": chapter,
	})
}

func UpdateChapter(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")

	chapter := domain.Chapter{}

	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get chapter data: " + err.Error(),
		})
		return
	}

	err := firestore_services.UpdateChapter(novelId, &chapter, ctx)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to update chapter: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update chapter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chapter": chapter,
	})
}

func DeleteChapter(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")
	chapterId := c.Param("chapter")

	err := firestore_services.DeleteChapter(novelId, chapterId, ctx)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to delete chapter: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete chapter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chapter deleted successfully",
	})
}
