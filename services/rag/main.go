package main

import (
	"adk"
	"log"
	"rag/methods"
	"rag/rag"
	"rag/settings"
)

func main() {
	ragSettings := settings.GetRAGSettings()

	apiKey := getOpenRouterAPIKey()
	if apiKey == "" {
		log.Fatalf("Failed to get OpenRouter API key from environment or vault.json")
	}

	embeddingClient := rag.NewEmbeddingClient(
		ragSettings.OpenRouterEmbeddingURL,
		ragSettings.OpenRouterEmbeddingModel,
		apiKey,
	)
	vectorStore := rag.NewVectorStore()
	indexer := rag.NewIndexer(embeddingClient, vectorStore)

	dbPath := rag.GetDatabasePath()
	log.Printf("Indexing database from: %s", dbPath)
	if err := indexer.IndexDatabase(dbPath); err != nil {
		log.Fatalf("Failed to index database: %v", err)
	}

	// Сохраняем компоненты в глобальном хранилище для использования в методах
	methods.InitializeRAG(embeddingClient, vectorStore, ragSettings.TopK)

	provider, err := adk.NewProvider("eliza")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Main server instance
	server, err := adk.NewServer(adk.ServerConfig{
		Port:               ":50057",
		Provider:           provider,
		SendMessageHandler: methods.SendMessage,
		GetTaskHandler:     methods.GetTask,
		Database:           "a2a",
		Collection:         "rag",
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Println("RAG server starting on :50057")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
