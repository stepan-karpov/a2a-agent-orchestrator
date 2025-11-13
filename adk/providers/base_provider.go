package providers

// Provider - interface for LLM providers
type Provider interface {
	ChatCompletion(messages []Message) (*ChatResponse, error)
}

// Message - message for provider
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse - response from provider
type ChatResponse struct {
	Content string
}
