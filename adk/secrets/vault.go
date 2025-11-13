package secrets

import (
	"encoding/json"
	"fmt"
	"os"
)

// Vault holds all application secrets
type Vault struct {
	ElizaToken       string `json:"eliza-token"`
	OpenRouterAPIKey string `json:"openrouter-api-key"`
	MongoDBURI       string `json:"mongodb-uri"`
}

// LoadVault loads secrets from vault.json file
func LoadVault(vaultPath string) (*Vault, error) {
	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	var vault Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, fmt.Errorf("failed to parse vault file: %w", err)
	}

	return &vault, nil
}
