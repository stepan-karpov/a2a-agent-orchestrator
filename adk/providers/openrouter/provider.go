package openrouter

import (
	agentsPkg "adk/agents"
	"adk/providers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Provider - implementation of OpenRouter provider
type Provider struct {
	Token string
}

// NewProvider creates a new OpenRouter provider with token
func NewProvider(token string) *Provider {
	return &Provider{
		Token: token,
	}
}

type openRouterRequest struct {
	Messages   []providers.Message `json:"messages"`
	Model      string              `json:"model"`
	Stream     bool                `json:"stream"`
	ToolChoice string              `json:"tool_choice,omitempty"`
	Tools      []openRouterTool    `json:"tools,omitempty"`
}

type openRouterTool struct {
	Type     string             `json:"type"`
	Function openRouterFunction `json:"function"`
}

type openRouterFunction struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Parameters  openRouterFunctionParameters `json:"parameters"`
}

type openRouterFunctionParameters struct {
	Type       string                        `json:"type"`
	Properties map[string]openRouterProperty `json:"properties"`
	Required   []string                      `json:"required"`
}

type openRouterProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type openRouterResponse struct {
	Choices []openRouterChoice `json:"choices"`
}

type openRouterChoice struct {
	Message openRouterMessage `json:"message"`
}

type openRouterMessage struct {
	Role      string               `json:"role"`
	Content   *string              `json:"content"` // Can be null
	ToolCalls []openRouterToolCall `json:"tool_calls,omitempty"`
}

type openRouterToolCall struct {
	ID       string                     `json:"id"`
	Type     string                     `json:"type"`
	Function openRouterToolCallFunction `json:"function"`
}

type openRouterToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

func (p *Provider) ChatCompletion(messages []providers.Message, agents []agentsPkg.Agent) (*providers.ChatResponse, error) {
	fmt.Println("OpenRouter messages request: ", messages)

	reqBody := openRouterRequest{
		Messages: messages,
		Model:    "deepseek/deepseek-chat-v3-0324",
		Stream:   false,
	}

	// Add tools if agents are provided
	if len(agents) > 0 {
		reqBody.ToolChoice = "auto"
		reqBody.Tools = make([]openRouterTool, 0, len(agents))

		for _, agent := range agents {
			tool := openRouterTool{
				Type: "function",
				Function: openRouterFunction{
					Name:        agent.Name,
					Description: agent.Description,
					Parameters: openRouterFunctionParameters{
						Type: "object",
						Properties: map[string]openRouterProperty{
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

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
	req.Header.Set("Content-Type", "application/json")

	// Log curl command for debugging
	fmt.Println("\n=== OpenRouter Request (curl) ===")
	fmt.Printf("curl -X POST https://openrouter.ai/api/v1/chat/completions \\\n")
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

	var openRouterResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openRouterResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := openRouterResp.Choices[0]
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

		fmt.Printf("OpenRouter tool call: agent=%s, message=%s\n", targetAgent.Name, agentQueryMessage)

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

	openRouterTextResponse := *message.Content
	fmt.Println("Response from OpenRouter: ", openRouterTextResponse)

	return &providers.ChatResponse{
		Content: openRouterTextResponse,
	}, nil
}
