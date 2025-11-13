package adk

import (
	"adk/providers"
	"adk/providers/openrouter"
	"adk/providers/eliza"
	"adk/secrets"
	"errors"
	"os"
	"path/filepath"
)

// NewProvider creates a provider by name
// It loads secrets from vault.json in the project root
func NewProvider(providerName string) (providers.Provider, error) {
	// Try to find vault.json in common locations
	var vault *secrets.Vault
	var err error

	// Try relative path from current working directory
	cwd, _ := os.Getwd()
	vaultPath := filepath.Join(cwd, "vault.json")
	vault, err = secrets.LoadVault(vaultPath)
	if err != nil {
		// Try from adk directory (go up two levels to project root)
		vaultPath = filepath.Join(cwd, "..", "..", "vault.json")
		vault, err = secrets.LoadVault(vaultPath)
		if err != nil {
			return nil, err
		}
	}

	switch providerName {
	case "eliza":
		if vault.ElizaToken == "" {
			return nil, errors.New("eliza-token not found in vault.json")
		}
		return eliza.NewProvider(vault.ElizaToken), nil
	case "openrouter":
		if vault.OpenRouterAPIKey == "" {
			return nil, errors.New("openrouter-api-key not found in vault.json")
		}
		return openrouter.NewProvider(vault.OpenRouterAPIKey), nil
	default:
		return nil, errors.New("unknown provider: " + providerName)
	}
}
