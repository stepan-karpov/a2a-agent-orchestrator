package methods

import (
	"adk"
	"adk/a2a/server"
	a2aServerProto "adk/a2a/server"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func GetTask(context context.Context, req *server.GetTaskRequest, server *adk.Server) (*server.Task, error) {
	fmt.Println("GetTaskRequest: ", req)

	task := &a2aServerProto.Task{
		Id:        uuid.New().String(),
		ContextId: "TODO: add context id",
		Status:    a2aServerProto.TaskState_TASK_STATE_COMPLETED,
		Artifacts: []*a2aServerProto.Artifact{
			{
				Type:    "text",
				Content: answer,
			},
		},
	}

	return task, nil
}
