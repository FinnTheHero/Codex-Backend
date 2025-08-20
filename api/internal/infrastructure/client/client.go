package firestore_client

import (
	cmn "Codex-Backend/api/internal/common"
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
}

func FirestoreClient() (*firestore.Client, error) {
	ctx := context.Background()

	credentials_json := cmn.GetEnvVariable("GOOGLE_CREDENTIALS")

	sa := option.WithCredentialsJSON([]byte(credentials_json))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Firebase App: " + err.Error()), Status: http.StatusInternalServerError}
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Firestore Client: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return client, nil
}
