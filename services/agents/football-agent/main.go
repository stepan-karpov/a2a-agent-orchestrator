package main

import (
	"adk"
	"football-agent/methods"
	"log"
)

func main() {
	provider, err := adk.NewProvider("eliza")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Main server instance
	server, err := adk.NewServer(adk.ServerConfig{
		Port:               ":50054",
		Provider:           provider,
		SendMessageHandler: methods.SendMessage,
		GetTaskHandler:     methods.GetTask,
		Database:           "a2a",
		Collection:         "football_agent",
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
