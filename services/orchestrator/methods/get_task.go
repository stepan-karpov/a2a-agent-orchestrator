package methods

import (
	"adk/a2a/server"
	"context"
	"log"
)

func GetTask(ctx context.Context, req *server.GetTaskRequest) (*server.Task, error) {
	log.Printf("Custom GetTask handler: %+v", req)

	task := &server.Task{
		Id:        req.Name,
		ContextId: "context-custom",
		Status:    server.TaskState_TASK_STATE_WORKING,
	}

	return task, nil
}
