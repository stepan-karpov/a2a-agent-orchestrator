package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"adk/execution"
	"context"
	"fmt"
	"log"
	"orchestrator/settings"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// callRAGService calls RAG service on localhost:50057 and returns the response
func callRAGService(ctx context.Context, query string, contextId string) (*a2aServerProto.SendMessageResponse, error) {
	url := "passthrough:///localhost:50057"

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RAG service: %w", err)
	}
	defer conn.Close()

	client := a2aServerProto.NewA2AServiceClient(conn)

	req := &a2aServerProto.SendMessageRequest{
		Request: &a2aServerProto.Message{
			ContextId: contextId,
			Role:      a2aServerProto.Role_ROLE_USER,
			Content:   query,
		},
	}

	if req.Request == nil {
		return nil, fmt.Errorf("failed to create message request: request is nil")
	}

	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to RAG service: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("invalid response from RAG service: response is nil")
	}

	return resp, nil
}

// extractDocumentsFromRAGResponse extracts documents from RAG service response
// RAG service returns completed task immediately in SendMessage response
func extractDocumentsFromRAGResponse(resp *a2aServerProto.SendMessageResponse) (string, error) {
	if resp == nil || resp.Task == nil {
		return "", fmt.Errorf("invalid response from RAG service: task is nil")
	}

	// RAG service returns completed task immediately
	if resp.Task.Status != a2aServerProto.TaskState_TASK_STATE_COMPLETED {
		return "", fmt.Errorf("RAG service returned task with status: %v", resp.Task.Status)
	}

	if len(resp.Task.Artifacts) == 0 {
		return "", fmt.Errorf("RAG service task completed but no artifacts found")
	}

	return resp.Task.Artifacts[0].Content, nil
}

func SendMessage(ctx context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	serverSettings := settings.GetServerSettings()

	// Call RAG service to get relevant documents
	ragDocuments := ""
	ragQuery := req.Request.Content
	if ragQuery != "" {
		ragResp, err := callRAGService(ctx, ragQuery, req.Request.ContextId)
		if err != nil {
			log.Printf("Warning: failed to get response from RAG service: %v", err)
			// Continue even if RAG service is unavailable
		} else {
			documents, err := extractDocumentsFromRAGResponse(ragResp)
			if err != nil {
				log.Printf("Warning: failed to extract documents from RAG response: %v", err)
				// Continue even if failed to extract documents
			} else {
				ragDocuments = documents
				log.Printf("Retrieved %d characters of documents from RAG service", len(ragDocuments))
			}
		}
	}

	// Build prompt with documents
	prompt := serverSettings.Prompt
	if ragDocuments != "" {
		prompt = strings.TrimSpace(prompt) + "\n\nRelevant documents from knowledge base:\n" + ragDocuments
	}

	// Convert ServerSettings to ExecutionSettings
	executionSettings := &execution.ExecutionSettings{
		Prompt:       prompt,
		HistoryLimit: serverSettings.HistoryLimit,
	}

	task, err := server.CreateNewDetachedTask(ctx, req.Request, executionSettings)
	if err != nil {
		log.Printf("Error creating new detached task: %v", err)
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("task is nil")
	}

	return &a2aServerProto.SendMessageResponse{Task: task}, nil
}
