package main

import (
	"adk"
	"log"
	"orchestrator/methods"
)

func main() {
	provider, err := adk.NewProvider("eliza")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Main server instance
	server, err := adk.NewServer(adk.ServerConfig{
		Port:               ":50051",
		Provider:           provider,
		SendMessageHandler: methods.SendMessage,
		GetTaskHandler:     methods.GetTask,
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
