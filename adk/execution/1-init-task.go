package execution

import (
	"adk/a2a/server"
	"context"
)

func InitTask(context context.Context, message *server.Message, taskId string) *server.Task {
	return &server.Task{
		Id:        taskId,
		ContextId: message.ContextId,
		Status:    server.TaskState_TASK_STATE_WORKING,
		History:   []*server.Message{message},
	}
}
