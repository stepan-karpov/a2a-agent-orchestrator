package eliza

import (
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
	Messages []providers.Message `json:"messages"`
	Model    string              `json:"model"`
}

type elizaResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *Provider) ChatCompletion(messages []providers.Message) (*providers.ChatResponse, error) {
	fmt.Println("Eliza messages request: ", messages)

	reqBody := elizaRequest{
		Messages: messages,
		Model:    "internal-deepseek",
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
	elizaTextResponse := elizaResp.Choices[0].Message.Content
	
	fmt.Println("Response from Eliza: ", elizaTextResponse)

	return &providers.ChatResponse{
		Content: elizaTextResponse,
	}, nil
}
