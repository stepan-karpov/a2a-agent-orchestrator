package mongo

import (
	"adk/a2a/server"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
)

// SaveTask saves a task to MongoDB as JSON
func SaveTask(ctx context.Context, task *server.Task, database, collection string) error {
	// Get MongoDB client
	client, err := GetClient(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Convert task to JSON first (to preserve protobuf field names and structures)
	jsonData, err := protojson.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task to JSON: %w", err)
	}

	// Convert JSON to BSON document
	var bsonDoc bson.M
	if err := json.Unmarshal(jsonData, &bsonDoc); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to BSON: %w", err)
	}

	// Add save timestamp
	bsonDoc["updated_at"] = time.Now()

	// Save to MongoDB
	coll := client.Database(database).Collection(collection)
	_, err = coll.ReplaceOne(
		ctx,
		bson.M{"id": task.Id},
		bsonDoc,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to save task to MongoDB: %w", err)
	}

	return nil
}
