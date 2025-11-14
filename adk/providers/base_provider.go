package providers

import "adk/agents"

// Provider - interface for LLM providers
type Provider interface {
	ChatCompletion(messages []Message, agents []agents.Agent) (*ChatResponse, error)
}

// Message - message for provider
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AgentQuery - represents a query to an agent when tool_calls are present
type AgentQuery struct {
	Agent   *agents.Agent // Agent to call
	Message string        // Message to send to the agent
}

// ChatResponse - response from provider
type ChatResponse struct {
	Content    string      // Regular response content
	AgentQuery *AgentQuery // Agent query when tool_calls are present (nil for regular responses)
}
