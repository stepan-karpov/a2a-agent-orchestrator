package deepseek

import "adk/providers"

// Provider - implementation of DeepSeek provider
type Provider struct{}

func (p *Provider) ChatCompletion(messages []providers.Message) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{
		Content: "HelloWorld from DeepSeek",
	}, nil
}
