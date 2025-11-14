package eliza

import (
	agentsPkg "adk/agents"
	"adk/providers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Provider - implementation of Eliza provider
type Provider struct {
	Token string
}

// NewProvider creates a new Eliza provider with token
func NewProvider(token string) *Provider {
	return &Provider{
		Token: token,
	}
}

type elizaRequest struct {
	Messages   []providers.Message `json:"messages"`
	Model      string              `json:"model"`
	Stream     bool                `json:"stream"`
	ToolChoice string              `json:"tool_choice,omitempty"`
	Tools      []elizaTool         `json:"tools,omitempty"`
}

type elizaTool struct {
	Type     string        `json:"type"`
	Function elizaFunction `json:"function"`
}

type elizaFunction struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Parameters  elizaFunctionParameters `json:"parameters"`
}

type elizaFunctionParameters struct {
	Type       string                   `json:"type"`
	Properties map[string]elizaProperty `json:"properties"`
	Required   []string                 `json:"required"`
}

type elizaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type elizaResponse struct {
	Choices []elizaChoice `json:"choices"`
}

type elizaChoice struct {
	Message elizaMessage `json:"message"`
}

type elizaMessage struct {
	Role      string          `json:"role"`
	Content   *string         `json:"content"` // Can be null
	ToolCalls []elizaToolCall `json:"tool_calls,omitempty"`
}

type elizaToolCall struct {
	ID       string                `json:"id"`
	Type     string                `json:"type"`
	Function elizaToolCallFunction `json:"function"`
}

type elizaToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

func (p *Provider) ChatCompletion(messages []providers.Message, agents []agentsPkg.Agent) (*providers.ChatResponse, error) {
	fmt.Println("Eliza messages request: ", messages)

	reqBody := elizaRequest{
		Messages: messages,
		Model:    "internal-deepseek",
		Stream:   false,
	}

	// Add tools if agents are provided
	if len(agents) > 0 {
		reqBody.ToolChoice = "auto"
		reqBody.Tools = make([]elizaTool, 0, len(agents))

		for _, agent := range agents {
			tool := elizaTool{
				Type: "function",
				Function: elizaFunction{
					Name:        agent.Name,
					Description: agent.Description,
					Parameters: elizaFunctionParameters{
						Type: "object",
						Properties: map[string]elizaProperty{
							"message": {
								Type:        "string",
								Description: "Сообщение для отправки удаленному A2A агенту",
							},
						},
						Required: []string{"message"},
					},
				},
			}
			reqBody.Tools = append(reqBody.Tools, tool)
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTP request will be successful with Yandex VPN only
	req, err := http.NewRequest("POST", "https://api.eliza.yandex.net/raw/auto/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Set("Content-Type", "application/json")

	// Log curl command for debugging
	fmt.Println("\n=== Eliza Request (curl) ===")
	fmt.Printf("curl -X POST https://api.eliza.yandex.net/raw/auto/chat/completions \\\n")
	fmt.Printf("  -H \"Authorization: Bearer %s\" \\\n", p.Token)
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d '%s'\n", string(jsonData))
	fmt.Println("==================================")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var elizaResp elizaResponse
	if err := json.NewDecoder(resp.Body).Decode(&elizaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(elizaResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := elizaResp.Choices[0]
	message := choice.Message

	// Check if response contains tool_calls
	if len(message.ToolCalls) > 0 {
		// Find the first tool call
		toolCall := message.ToolCalls[0]

		// Find the agent by function name
		var targetAgent *agentsPkg.Agent
		for i := range agents {
			if agents[i].Name == toolCall.Function.Name {
				agentCopy := agents[i] // Copy to avoid taking address of loop variable
				targetAgent = &agentCopy
				break
			}
		}

		if targetAgent == nil {
			return nil, fmt.Errorf("agent not found for tool call: %s", toolCall.Function.Name)
		}

		// Parse arguments to extract message
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			return nil, fmt.Errorf("failed to parse tool call arguments: %w", err)
		}

		// Extract message from arguments
		agentQueryMessage, ok := args["message"].(string)
		if !ok {
			return nil, fmt.Errorf("message not found in tool call arguments")
		}

		fmt.Printf("Eliza tool call: agent=%s, message=%s\n", targetAgent.Name, agentQueryMessage)

		return &providers.ChatResponse{
			AgentQuery: &providers.AgentQuery{
				Agent:   targetAgent,
				Message: agentQueryMessage,
			},
		}, nil
	}

	// Regular response with content
	if message.Content == nil {
		return nil, fmt.Errorf("response has no content and no tool_calls")
	}

	elizaTextResponse := *message.Content
	fmt.Println("Response from Eliza: ", elizaTextResponse)

	return &providers.ChatResponse{
		Content: elizaTextResponse,
	}, nil
}
