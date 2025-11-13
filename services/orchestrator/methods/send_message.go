package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"adk/execution"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func getServerSettings() *execution.ExecutionSettings {
	// Try multiple paths for orchestrator_setting.json
	paths := []string{
		"orchestrator_setting.json",
		"services/orchestrator/orchestrator_setting.json",
		filepath.Join("..", "orchestrator_setting.json"),
	}

	var serverSettingsJson []byte
	var err error
	for _, path := range paths {
		serverSettingsJson, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("Error reading orchestrator setting file: %v", err)
	}

	var serverSettings execution.ExecutionSettings
	if err := json.Unmarshal(serverSettingsJson, &serverSettings); err != nil {
		log.Fatalf("Error parsing orchestrator setting file: %v", err)
	}

	return &serverSettings
}

func SendMessage(context context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	serverSetting := getServerSettings()
	task, err := server.CreateNewDetachedTask(context, req.Request, serverSetting)
	if err != nil {
		log.Printf("Error creating new detached task: %v", err)
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("task is nil")
	}

	return &a2aServerProto.SendMessageResponse{Task: task}, nil
}
