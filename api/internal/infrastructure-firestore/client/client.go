package firestore_client

import (
	"Codex-Backend/api/internal/config"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
}

func NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	creds, err := config.GetEnvVariable("FIRESTORE_JSON")
	if err != nil {
		return nil, err
	}

	projectID, err := config.GetEnvVariable("PROJECT_ID")
	if err != nil {
		return nil, err
	}

	return firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(creds)))
}
