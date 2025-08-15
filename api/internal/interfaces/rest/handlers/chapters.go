package firestore_handlers

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_services "Codex-Backend/api/internal/usecases/collections"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func GetPaginatedChapters(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")
	if novelId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel ID not found",
		})
		return
	}

	options := domain.CursorOptions{
		NovelID: novelId,
		Cursor:  0,
		Limit:   100,
		SortBy:  firestore.Desc,
	}

	if cursor, exists := c.GetQuery("cursor"); exists {
		curs, err := strconv.Atoi(cursor)
		if err == nil {
			options.Cursor = curs
		}
	}

	if limit, exists := c.GetQuery("limit"); exists {
		lim, err := strconv.Atoi(limit)
		if err == nil {
			options.Limit = lim
		}
	}

	if sortBy, exists := c.GetQuery("sort"); exists {
		switch sortBy {
		case "asc":
			options.SortBy = firestore.Asc
		case "desc":
			options.SortBy = firestore.Desc
		default:
			options.SortBy = firestore.Desc
		}
	}

	response, err := firestore_services.GetCursorPaginatedChapters(options, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to retrieve chapters: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve chapters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chapters":    &response.Chapters,
		"next_cursor": &response.NextCursor,
	})
}

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
	if e, ok := err.(*cmn.Error); ok {
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
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to retrieve chapters: " + e.Error(),
		})
		return
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

func BatchUploadChapters(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	novelId := c.Param("novel")

	if novelId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Novel ID not found",
		})
		return
	}

	var chapters []domain.Chapter
	if err := c.ShouldBindJSON(&chapters); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get chapters data: " + err.Error(),
		})
		return
	}

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Nothing to upload",
		})
		return
	}

	err := firestore_services.BatchUploadChapters(novelId, chapters, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Failed to upload chapters: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload chapters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chapters uploaded successfully",
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
	if e, ok := err.(*cmn.Error); ok {
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
	if e, ok := err.(*cmn.Error); ok {
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
	if e, ok := err.(*cmn.Error); ok {
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
