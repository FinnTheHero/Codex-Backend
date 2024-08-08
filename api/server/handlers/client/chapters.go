package client_handler

import (
	aws_methods "Codex-Backend/api/aws/methods"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindChapter(c *gin.Context) {
	novel := c.Param("novel")
	chapter := c.Param("chapter")

	tables, err := aws_methods.GetTables()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get tables"})
		return
	}

	for _, table := range tables {
		if table == novel {
			ch, err := aws_methods.GetChapter(novel, chapter)
			if err != nil {
				c.JSON(http.StatusInternalServerError,
					gin.H{"error": err.Error()},
				)
				return
			}

			c.JSON(http.StatusFound,
				gin.H{"chapter": ch},
			)
			return
		}
	}

}

func FindAllChapters(c *gin.Context) {
	novel := c.Param("novel")

	tables, err := aws_methods.GetTables()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, table := range tables {
		if table == novel {
			chapters, err := aws_methods.GetAllChapters(novel)
			if err != nil {
				c.JSON(http.StatusInternalServerError,
					gin.H{"error": err.Error()},
				)
				return
			}

			c.JSON(http.StatusFound,
				gin.H{"chapters": chapters},
			)
			return
		}
	}

	c.JSON(http.StatusNotFound,
		gin.H{"error": "Novel not found"},
	)
}
