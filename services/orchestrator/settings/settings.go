package settings

import (
	"adk/agents"
	"adk/execution"
	"encoding/json"
	"log"
	"os"
)

// ServerSettings extends ExecutionSettings with agents support
type ServerSettings struct {
	execution.ExecutionSettings
	Agents []agents.Agent `json:"agents"`
}

func GetServerSettings() *ServerSettings {
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

	var serverSettings ServerSettings
	if err := json.Unmarshal(serverSettingsJson, &serverSettings); err != nil {
		log.Fatalf("Error parsing orchestrator setting file: %v", err)
	}

	// Initialize agents slice if nil
	if serverSettings.Agents == nil {
		serverSettings.Agents = make([]agents.Agent, 0)
	}

	return &serverSettings
}
