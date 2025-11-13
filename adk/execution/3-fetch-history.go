package execution

import (
	"adk/a2a/server"
	"adk/providers"
	"adk/storages/mongo"
	"context"
	"fmt"
)

func FetchHistory(ctx context.Context, task *server.Task, database, collection string) ([]providers.Message, error) {
	// Fetch all tasks by context_id, sorted from oldest to newest
	tasks, err := mongo.FetchTasks(ctx, task.ContextId, database, collection)
	fmt.Println("Number of tasks fetched from history: ", len(tasks))
	if err != nil {
		return nil, err
	}

	var history []providers.Message

	// Iterate through tasks and extract messages
	for _, t := range tasks {
		// Extract user message from metadata
		if t.Metadata != nil && t.Metadata.Fields != nil {
			if messageValue, ok := t.Metadata.Fields["message"]; ok {
				if messageStr := messageValue.GetStringValue(); messageStr != "" {
					history = append(history, providers.Message{
						Role:    "user",
						Content: messageStr,
					})
				}
			}
		}

		// Extract assistant responses from artifacts
		for _, artifact := range t.Artifacts {
			if artifact.Type == "text" && artifact.Content != "" {
				history = append(history, providers.Message{
					Role:    "assistant",
					Content: artifact.Content,
				})
			}
		}
	}

	return history, nil
}
