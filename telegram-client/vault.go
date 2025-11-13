package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetTelegramToken() string {
	// First try to get from environment variable
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		return token
	}

	// Fallback to vault.json
	// Try different paths
	paths := []string{
		"../vault.json",
		"/app/vault.json",
		"vault.json",
	}

	cwd, _ := os.Getwd()
	paths = append(paths, filepath.Join(cwd, "..", "vault.json"))

	for _, vaultPath := range paths {
		data, err := os.ReadFile(vaultPath)
		if err != nil {
			continue
		}
		var vaultData map[string]interface{}
		if err := json.Unmarshal(data, &vaultData); err != nil {
			log.Printf("Failed to parse vault.json: %v", err)
			continue
		}
		if token, ok := vaultData["telegram-bot-token"].(string); ok && token != "" {
			return token
		}
	}

	log.Fatal("Failed to get Telegram bot token from environment or vault.json")
	return ""
}
