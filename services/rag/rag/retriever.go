package rag

import (
	"fmt"
)

// Retriever performs a search for relevant documents
type Retriever struct {
	embeddingClient *EmbeddingClient
	vectorStore     *VectorStore
	topK            int
}

// NewRetriever creates a new retriever
func NewRetriever(embeddingClient *EmbeddingClient, vectorStore *VectorStore, topK int) *Retriever {
	return &Retriever{
		embeddingClient: embeddingClient,
		vectorStore:     vectorStore,
		topK:            topK,
	}
}

// Retrieve finds the TopK most relevant documents for the query
func (r *Retriever) Retrieve(query string) ([]string, []float32, error) {
	// Generate an embedding for the query
	queryEmbedding, err := r.embeddingClient.GetEmbedding(query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search for relevant documents
	documents, scores, err := r.vectorStore.Search(queryEmbedding, r.topK)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search documents: %w", err)
	}

	return documents, scores, nil
}
