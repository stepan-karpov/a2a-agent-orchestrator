package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"
	"fmt"
	"strings"
	"sync"

	"rag/rag"

	"github.com/google/uuid"
)

var (
	ragRetriever *rag.Retriever
	ragInitOnce  sync.Once
)

func InitializeRAG(embeddingClient *rag.EmbeddingClient, vectorStore *rag.VectorStore, topK int) {
	ragInitOnce.Do(func() {
		ragRetriever = rag.NewRetriever(embeddingClient, vectorStore, topK)
	})
}

func SendMessage(ctx context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	query := req.Request.Content
	fmt.Printf("RAG - Query: %s\n", query)

	if ragRetriever == nil {
		return nil, fmt.Errorf("RAG system not initialized")
	}

	// Ищем релевантные документы
	documents, scores, err := ragRetriever.Retrieve(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %w", err)
	}

	// Формируем ответ с найденными документами
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d relevant documents:\n\n", len(documents)))

	for i, doc := range documents {
		result.WriteString(fmt.Sprintf("%d. (Score: %.3f) %s\n\n", i+1, scores[i], doc))
	}

	// Если документов не найдено
	if len(documents) == 0 {
		result.WriteString("No relevant documents found for your query.")
	}

	task := &a2aServerProto.Task{
		Id:        uuid.New().String(),
		ContextId: req.Request.ContextId,
		Status:    a2aServerProto.TaskState_TASK_STATE_COMPLETED,
		Artifacts: []*a2aServerProto.Artifact{
			{
				Type:    "text",
				Content: result.String(),
			},
		},
	}

	return &a2aServerProto.SendMessageResponse{Task: task}, nil
}
