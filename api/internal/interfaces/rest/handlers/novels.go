package firestore_handlers

import (
	cmn "Codex-Backend/api/internal/common"
	queue "Codex-Backend/api/internal/common/river"
	"Codex-Backend/api/internal/domain"
	firestore_services "Codex-Backend/api/internal/usecases/collections"
	"Codex-Backend/api/internal/usecases/worker"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func EPUBNovel(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	defer func() {
		if c.Request.MultipartForm != nil {
			c.Request.MultipartForm.RemoveAll()
		}
	}()

	epubFile, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get EPUB file: " + err.Error(),
		})
		return
	}

	maxSize := int64(32 * 1024 * 1024) // 32MB limit
	if epubFile.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 32MB)"})
		return
	}

	file, err := epubFile.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to open EPUB file: " + err.Error(),
		})
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	riverClient := queue.GetRiverClient(ctx)
	_, err = riverClient.Insert(ctx, worker.ProcessEPUBArgs{File: fileData}, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "EPUB file uploaded successfully",
	})
}

func FindNovel(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	param := c.Param("novel")

	withId := false
	withTitle := false

	if strings.HasPrefix(param, "novel_") {
		withId = true
	} else {
		withTitle = true
	}

	if withId {
		novel, err := firestore_services.GetNovelById(param, ctx)
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
	} else if withTitle {
		novel, err := firestore_services.GetNovelByTitle(param, ctx)
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
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel Title and ID not found",
		})
		return
	}
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

	err, id := firestore_services.CreateNovel(novel, ctx)
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
		"id":      id,
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
