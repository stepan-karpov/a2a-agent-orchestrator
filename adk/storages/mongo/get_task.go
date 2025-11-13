package mongo

import (
	"adk/a2a/server"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/encoding/protojson"
)

// GetTask retrieves a task from MongoDB by ID
func GetTask(ctx context.Context, taskId string, database, collection string) (*server.Task, error) {
	// Get MongoDB client
	client, err := GetClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	// Get collection
	coll := client.Database(database).Collection(collection)

	// Find task by ID
	var bsonDoc bson.M
	err = coll.FindOne(ctx, bson.M{"id": taskId}).Decode(&bsonDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("task not found: %s", taskId)
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	// Remove MongoDB-specific fields
	delete(bsonDoc, "updated_at")
	delete(bsonDoc, "_id")

	// Convert BSON to JSON, then to Task (to handle structpb.Struct properly)
	jsonData, err := bson.MarshalExtJSON(bsonDoc, false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal BSON to JSON: %w", err)
	}

	var task server.Task
	if err := protojson.Unmarshal(jsonData, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to Task: %w", err)
	}

	return &task, nil
}
