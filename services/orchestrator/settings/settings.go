package settings

import (
	"adk/execution"
	"encoding/json"
	"log"
	"os"
)

func GetServerSettings() *execution.ExecutionSettings {
	// Try multiple paths for orchestrator_setting.json
	paths := []string{
		"orchestrator_setting.json",
		"services/orchestrator/orchestrator_setting.json",
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
