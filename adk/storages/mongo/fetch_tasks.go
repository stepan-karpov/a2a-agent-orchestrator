package mongo

import (
	"adk/a2a/server"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
)

// FetchTasks retrieves all tasks by context_id, sorted from oldest to newest
func FetchTasks(ctx context.Context, contextId string, database, collection string) ([]*server.Task, error) {
	// Get MongoDB client
	client, err := GetClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	// Get collection
	coll := client.Database(database).Collection(collection)

	// Find all tasks by context_id, sorted by updated_at ascending (oldest first)
	findOptions := options.Find().SetSort(bson.M{"updated_at": 1})
	cursor, err := coll.Find(ctx, bson.M{"contextId": contextId}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*server.Task

	// Iterate through results
	for cursor.Next(ctx) {
		var bsonDoc bson.M
		if err := cursor.Decode(&bsonDoc); err != nil {
			return nil, fmt.Errorf("failed to decode task: %w", err)
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

		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return tasks, nil
}
