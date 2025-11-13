package eliza

import "adk/providers"

// Provider - implementation of Eliza provider
type Provider struct{}

func (p *Provider) ChatCompletion(messages []providers.Message) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{
		Content: "HelloWorld from Eliza",
	}, nil
}
