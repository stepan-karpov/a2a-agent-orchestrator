package deepseek

import "adk/providers"

// Provider - реализация провайдера Eliza
type Provider struct{}

// ChatCompletion реализует метод интерфейса providers.Provider
func (p *Provider) ChatCompletion(messages []providers.Message) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{
		Content: "HelloWorld from DeepSeek",
	}, nil
}
