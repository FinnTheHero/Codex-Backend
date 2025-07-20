package firestore_handlers

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_services "Codex-Backend/api/internal/usecases/collections"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindNovel(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")
	if novelId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel ID not found",
		})
		return
	}

	novel, err := firestore_services.GetNovel(novelId, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to retrieve novel: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve novel: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novel": novel,
	})
}

func FindAllNovels(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novels, err := firestore_services.GetAllNovels(ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to retrieve novels: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve novels: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novels": novels,
	})
}

func CreateNovel(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novel := domain.Novel{}

	if err := c.ShouldBindJSON(&novel); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get novel data: " + err.Error(),
		})
		return
	}

	err := firestore_services.CreateNovel(novel, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to create novel: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create novel: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Novel created successfully",
	})
}

func UpdateNovel(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")
	if novelId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel ID not found",
		})
		return
	}

	novel := domain.Novel{}

	if err := c.ShouldBindJSON(&novel); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get novel data: " + err.Error(),
		})
		return
	}

	err := firestore_services.UpdateNovel(novelId, novel, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to update novel: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update novel: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Novel updated successfully",
	})
}

func DeleteNovel(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")
	if novelId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel ID not found",
		})
		return
	}

	err := firestore_services.DeleteNovel(novelId, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to delete novel: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete novel: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Novel deleted successfully",
	})
}
