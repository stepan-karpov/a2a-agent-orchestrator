package rag

import (
	"math"
	"sync"
)

// Document represents a document with its embedding
type Document struct {
	Text      string
	Embedding []float32
}

// VectorStore stores documents and their embeddings
type VectorStore struct {
	documents []Document
	mu        sync.RWMutex
}

// NewVectorStore creates a new vector store
func NewVectorStore() *VectorStore {
	return &VectorStore{
		documents: make([]Document, 0),
	}
}

// Add adds a document to the store
func (vs *VectorStore) Add(text string, embedding []float32) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	vs.documents = append(vs.documents, Document{
		Text:      text,
		Embedding: embedding,
	})
}

// Search finds the TopK most relevant documents
func (vs *VectorStore) Search(queryEmbedding []float32, topK int) ([]string, []float32, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	if len(vs.documents) == 0 {
		return []string{}, []float32{}, nil
	}

	type scoredDoc struct {
		text  string
		score float32
	}

	scores := make([]scoredDoc, 0, len(vs.documents))

	for _, doc := range vs.documents {
		score := cosineSimilarity(queryEmbedding, doc.Embedding)
		scores = append(scores, scoredDoc{
			text:  doc.Text,
			score: score,
		})
	}

	// Sort by descending score
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[i].score < scores[j].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Return TopK
	if topK > len(scores) {
		topK = len(scores)
	}

	results := make([]string, 0, topK)
	resultScores := make([]float32, 0, topK)

	for i := 0; i < topK; i++ {
		results = append(results, scores[i].text)
		resultScores = append(resultScores, scores[i].score)
	}

	return results, resultScores, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
