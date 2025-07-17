package firestore_services

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	f_client "Codex-Backend/api/internal/infrastructure-firestore/client"
	f_collections "Codex-Backend/api/internal/infrastructure-firestore/collections"
	"context"
	"errors"
	"net/http"
)

func LoginUser(credentials domain.Credentials, ctx context.Context) (string, *domain.User, error) {
	client, err := f_client.FirestoreClient()
	if err != nil {
		return "", nil, err
	}
	defer client.Close()

	c := f_collections.Client{Client: client}

	user, err := c.GetUserByEmail(credentials.Email, ctx)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, &cmn.Error{Err: errors.New("Login Service Error - User not found"), Status: http.StatusNotFound}
	}

	err = VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return "", nil, &cmn.Error{Err: errors.New("Login Service Error - Invalid password"), Status: http.StatusUnauthorized}
	}

	token, err := GenerateToken(credentials.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
