package execution

import (
	"adk/a2a/server"
	"context"

	"google.golang.org/protobuf/types/known/structpb"
)

func InitTask(context context.Context, message *server.Message, taskId string) *server.Task {
	return &server.Task{
		Id:        taskId,
		ContextId: message.ContextId,
		Status:    server.TaskState_TASK_STATE_WORKING,
		Metadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"message": structpb.NewStringValue(message.Content),
			},
		},
	}
}
