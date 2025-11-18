package main

import (
	"adk/secrets"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// getOpenRouterAPIKey получает API ключ OpenRouter из переменной окружения или vault.json
func getOpenRouterAPIKey() string {
	paths := []string{
		"vault.json",
		"../vault.json",
		"../../vault.json",
		"../../../vault.json",
	}

	cwd, _ := os.Getwd()
	paths = append(paths,
		filepath.Join(cwd, "vault.json"),
		filepath.Join(cwd, "..", "vault.json"),
		filepath.Join(cwd, "..", "..", "vault.json"),
		filepath.Join(cwd, "..", "..", "..", "vault.json"),
	)

	for _, vaultPath := range paths {
		vault, err := secrets.LoadVault(vaultPath)
		if err == nil && vault.OpenRouterAPIKey != "" {
			return vault.OpenRouterAPIKey
		}
	}

	// Fallback: пробуем прочитать напрямую
	for _, vaultPath := range paths {
		data, err := os.ReadFile(vaultPath)
		if err != nil {
			continue
		}

		var vaultData map[string]interface{}
		if err := json.Unmarshal(data, &vaultData); err != nil {
			continue
		}

		if key, ok := vaultData["openrouter-api-key"].(string); ok && key != "" {
			return key
		}
	}

	log.Printf("Warning: OpenRouter API key not found in environment or vault.json")
	return ""
}
