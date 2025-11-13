package adk

import (
	"adk/providers"
	"adk/providers/deepseek"
	"adk/providers/eliza"
	"errors"
)

// NewProvider creates a provider by name
func NewProvider(providerName string) (providers.Provider, error) {
	switch providerName {
	case "eliza":
		return &eliza.Provider{}, nil
	case "deepseek":
		return &deepseek.Provider{}, nil
	default:
		return nil, errors.New("unknown provider: " + providerName)
	}
}
