package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"adk/execution"
	"context"
	"fmt"
	"log"
	"orchestrator/settings"
)

func SendMessage(context context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	serverSettings := settings.GetServerSettings()

	// Convert ServerSettings to ExecutionSettings
	executionSettings := &execution.ExecutionSettings{
		Prompt:       serverSettings.Prompt,
		HistoryLimit: serverSettings.HistoryLimit,
	}

	task, err := server.CreateNewDetachedTask(context, req.Request, executionSettings)
	if err != nil {
		log.Printf("Error creating new detached task: %v", err)
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("task is nil")
	}

	return &a2aServerProto.SendMessageResponse{Task: task}, nil
}
