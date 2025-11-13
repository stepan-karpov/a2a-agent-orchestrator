package methods

import (
	"adk"
	"adk/a2a/server"
	"context"
)

func GetTask(context context.Context, req *server.GetTaskRequest, server *adk.Server) (*server.Task, error) {
	task, err := server.GetTask(context, req)
	if err != nil {
		return nil, err
	}
	return task, nil
}
