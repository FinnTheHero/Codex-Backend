package firestore_services

import (
	"Codex-Backend/api/internal/common"
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure-firestore/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure-firestore/collections"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginUser(credentials domain.Credentials, ctx context.Context) (string, *domain.User, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return "", nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	user, err := c.GetUserByEmail(credentials.Email, ctx)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, &cmn.Error{Err: errors.New("Login Service Error - User not found"), Status: http.StatusNotFound}
	}

	err = cmn.VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return "", nil, &cmn.Error{Err: errors.New("Login Service Error - Invalid password"), Status: http.StatusUnauthorized}
	}

	token, err := cmn.GenerateToken(credentials.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func RegisterUser(newUser domain.NewUser, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	user, err := c.GetUserByEmail(newUser.Email, ctx)
	if e, ok := err.(*common.Error); ok {
		if e.StatusCode() != http.StatusNotFound {
			return &cmn.Error{Err: errors.New("Register Service Error - Getting User By Email: " + err.Error()), Status: http.StatusInternalServerError}
		}
	}

	if user != nil {
		return &cmn.Error{Err: errors.New("Register Service Error - User With Email " + newUser.Email + " Already Exists"), Status: http.StatusConflict}
	}

	id, err := cmn.GenerateID("user")
	if err != nil {
		return err
	}

	hashedPassword, err := cmn.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	err = c.CreateUser(domain.User{
		ID:        id,
		Username:  newUser.Username,
		Password:  string(hashedPassword),
		Email:     newUser.Email,
		Type:      "User",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, ctx)
	if err != nil {
		return err
	}

	return nil
}

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
