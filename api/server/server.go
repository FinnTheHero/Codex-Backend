package server

import (
	aws_methods "Codex-Backend/api/aws/methods"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server(svc *dynamodb.DynamoDB) {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	r.Use(cors.Default())

	// Find All Novels
	r.GET("/novels", func(c *gin.Context) {
		// TODO
	})

	// Find Novel
	r.GET("/:novel", func(c *gin.Context) {
		title := c.Param("novel")

		NovelSchema, err := aws_methods.GetNovel(svc, title)
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"novel": NovelSchema.Novel,
		})
	})

	r.Run()
}
