package mongo

import (
	"adk/secrets"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClient returns a MongoDB client connected to the database
func GetClient(ctx context.Context) (*mongo.Client, error) {
	// Load MongoDB URI from vault
	vaultPath := filepath.Join("..", "..", "vault.json")
	vault, err := secrets.LoadVault(vaultPath)
	if err != nil {
		cwd, _ := os.Getwd()
		vaultPath = filepath.Join(cwd, "vault.json")
		vault, err = secrets.LoadVault(vaultPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load vault: %w", err)
		}
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(vault.MongoDBURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return client, nil
}
