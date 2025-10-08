package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
)

// Client is a client for interacting with Firebase.
type Client struct {
	fs *firestore.Client
}

// NewClient creates a new Firebase client.
func NewClient(ctx context.Context, projectID string) (*Client, error) {
	fs, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Client{fs: fs}, nil
}
