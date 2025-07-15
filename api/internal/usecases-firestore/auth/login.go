package firestore_services

import (
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure-firestore/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure-firestore/collections"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoginUser(credentials domain.Credentials, ctx context.Context) (string, *domain.User, error) {
	client, err := firestore_client.NewFirestoreClient(ctx)
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
		return "", nil, status.Errorf(codes.Unauthenticated, "User not found")
	}

	err = VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return "", nil, status.Errorf(codes.Unauthenticated, "Invalid password")
	}

	token, err := GenerateToken(credentials.Email)
	if err != nil {
		return "", nil, status.Errorf(codes.Internal, "Error generating token: %v", err)
	}

	return token, user, nil
}
