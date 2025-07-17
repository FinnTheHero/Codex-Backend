package firestore_services

import (
	cmn "Codex-Backend/api/internal/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutUser(c *gin.Context) error {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		return &cmn.Error{Err: errors.New("Logout Service Error - Getting Cookie: " + err.Error()), Status: http.StatusBadRequest}
	}

	if tokenString == "" {
		return &cmn.Error{Err: errors.New("Logout Service Error - Token not found"), Status: http.StatusBadRequest}
	}

	c.SetCookie("Authorization", "", -1, "", "", true, true)

	return nil
}
