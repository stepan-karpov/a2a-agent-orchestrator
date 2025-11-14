package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"

	"github.com/google/uuid"
)

const (
	answer = "Here is an information about the author: John Doe. He is a software engineer at Google. He is 30 years old. He is from the United States."
)

func SendMessage(ctx context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	task := &a2aServerProto.Task{
		Id:        uuid.New().String(),
		ContextId: req.Request.ContextId,
		Status:    a2aServerProto.TaskState_TASK_STATE_COMPLETED,
		Artifacts: []*a2aServerProto.Artifact{
			{
				Type:    "text",
				Content: answer,
			},
		},
	}

	return &a2aServerProto.SendMessageResponse{Task: task}, nil
}
