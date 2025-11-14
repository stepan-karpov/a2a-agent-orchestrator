package main

import (
	"adk"
	"adk/agents"
	"log"
	"orchestrator/methods"
	"orchestrator/settings"
)

func main() {
	provider, err := adk.NewProvider("openrouter")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Main server instance
	server, err := adk.NewServer(adk.ServerConfig{
		Port:               ":50051",
		Provider:           provider,
		SendMessageHandler: methods.SendMessage,
		GetTaskHandler:     methods.GetTask,
		Database:           "a2a",
		Collection:         "orchestrator",
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	serverSettings := settings.GetServerSettings()

	for _, agent := range serverSettings.Agents {
		server.RegisterNewAgent(agents.Agent{
			Name:        agent.Name,
			Description: agent.Description,
			Url:         agent.Url,
		})
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
