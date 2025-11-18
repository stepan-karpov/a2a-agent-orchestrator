package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EmbeddingClient struct {
	apiURL         string
	apiKey         string
	model          string
	embeddingModel string
}

func NewEmbeddingClient(embeddingURL, embeddingModel, apiKey string) *EmbeddingClient {
	return &EmbeddingClient{
		apiURL:         embeddingURL,
		apiKey:         apiKey,
		model:          "openrouter",
		embeddingModel: embeddingModel,
	}
}

// GetEmbedding generates an embedding for the text
func (c *EmbeddingClient) GetEmbedding(text string) ([]float32, error) {
	return c.getOpenRouterEmbedding(text)
}

// getOpenRouterEmbedding gets an embedding through the OpenRouter API
func (c *EmbeddingClient) getOpenRouterEmbedding(text string) ([]float32, error) {
	modelName := c.embeddingModel
	if modelName == "" {
		modelName = "openai/text-embedding-3-small"
	}

	reqBody := map[string]interface{}{
		"model": modelName,
		"input": text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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

	var openRouterResp struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openRouterResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding in response")
	}

	return openRouterResp.Data[0].Embedding, nil
}
