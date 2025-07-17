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
)

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

	id, err := GenerateID("user")
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(newUser.Password)
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
