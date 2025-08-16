package firestore_server

import (
	"github.com/gin-gonic/gin"
)

func Server() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()

	RegisteredRoutes(r)

	r.Run()
}
