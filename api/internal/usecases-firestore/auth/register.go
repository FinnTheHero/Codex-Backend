package firestore_services

import (
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure-firestore/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure-firestore/collections"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterUser(newUser domain.NewUser, ctx context.Context) error {
	client, err := firestore_client.NewFirestoreClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	u, err := c.GetUserByEmail(newUser.Email, ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "Error searching for user with email: %v", err)
	}

	if u != nil {
		return status.Errorf(codes.AlreadyExists, "User with email %s already exists", newUser.Email)
	}

	id, err := GenerateID("user")
	if err != nil {
		return status.Errorf(codes.Internal, "Error generating ID: %v", err)
	}

	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		return status.Errorf(codes.Internal, "Error hashing password: %v", err)
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
