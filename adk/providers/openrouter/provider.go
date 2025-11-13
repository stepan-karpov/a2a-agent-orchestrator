package openrouter

import (
	"adk/providers"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Provider - implementation of DeepSeek provider via OpenRouter
type Provider struct {
	Token string
}

// NewProvider creates a new DeepSeek provider with token
func NewProvider(token string) *Provider {
	return &Provider{
		Token: token,
	}
}

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []providers.Message `json:"messages"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *Provider) ChatCompletion(messages []providers.Message) (*providers.ChatResponse, error) {
	reqBody := openRouterRequest{
		Model:    "deepseek/deepseek-chat-v3-0324",
		Messages: messages,
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
	deepSeekTextResponse := openRouterResp.Choices[0].Message.Content

	fmt.Println("Response from DeepSeek: ", deepSeekTextResponse)

	return &providers.ChatResponse{
		Content: deepSeekTextResponse,
	}, nil
}
