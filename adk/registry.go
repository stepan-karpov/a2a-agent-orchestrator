package adk

import (
	"adk/providers"
	"adk/providers/deepseek"
	"adk/providers/eliza"
	"errors"
	"fmt"
)

type Entry struct {
	ID       string
	Name     string
	Provider providers.Provider
}

// NewProvider создает провайдер по имени
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

func CreateEntry(providerName string) (Entry, error) {
	provider, err := NewProvider(providerName)
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		ID:       "1",
		Name:     "Test",
		Provider: provider,
	}, nil
}

func (e *Entry) PrintHello() {
	fmt.Println("Hello, World!")
}
