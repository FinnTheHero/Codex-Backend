package firestore_session

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
	return firestore.NewClient(ctx, "YOUR_PROJECT_ID", option.WithCredentialsJSON([]byte(creds)))
}
