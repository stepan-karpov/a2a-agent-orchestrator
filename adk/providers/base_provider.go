package providers

type Provider interface {
	ChatCompletion(messages []Message) (*ChatResponse, error)
}

type Message struct {
	Role    string
	Content string
}

type ChatResponse struct {
	Content string
}
