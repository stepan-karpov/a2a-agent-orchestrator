package settings

import (
	"encoding/json"
	"log"
	"os"
)

type RAGSettings struct {
	TopK                     int    `json:"top_k"`
	EmbeddingDimension       int    `json:"embedding_dimension"`
	OpenRouterEmbeddingURL   string `json:"openrouter_embedding_url"`
	OpenRouterEmbeddingModel string `json:"openrouter_embedding_model"`
}

func GetRAGSettings() *RAGSettings {
	paths := []string{
		"settings.json",
		"services/rag/settings.json",
		"./settings.json",
	}

	var settingsJson []byte
	var err error
	for _, path := range paths {
		settingsJson, err = os.ReadFile(path)
		if err == nil {
			log.Printf("Loaded settings from: %s", path)
			break
		}
	}
	if err != nil {
		log.Fatalf("Error reading settings file: %v", err)
	}

	var settings RAGSettings
	if err := json.Unmarshal(settingsJson, &settings); err != nil {
		log.Fatalf("Error parsing settings file: %v", err)
	}

	return &settings
}
