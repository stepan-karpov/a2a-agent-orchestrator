package methods

import (
	"adk/a2a/server"
	"adk/providers"
	"context"
	"log"
)

func SendMessage(ctx context.Context, req *server.SendMessageRequest, provider providers.Provider) (*server.SendMessageResponse, error) {
	log.Printf("Custom SendMessage handler: %+v", req)

	messages := []providers.Message{
		{
			Role:    req.Request.Role.String(),
			Content: req.Request.Content,
		},
	}

	response, err := provider.ChatCompletion(messages)
	if err != nil {
		log.Printf("Error calling provider: %v", err)
		return nil, err
	}

	responseMessage := &server.Message{
		Role:    server.Role_ROLE_ASSISTANT,
		Content: response.Content,
	}

	task := &server.Task{
		Id:        "task-custom",
		ContextId: req.Request.ContextId,
		Status:    server.TaskState_TASK_STATE_COMPLETED,
		History:   []*server.Message{req.Request, responseMessage},
	}

	return &server.SendMessageResponse{Task: task}, nil
}
