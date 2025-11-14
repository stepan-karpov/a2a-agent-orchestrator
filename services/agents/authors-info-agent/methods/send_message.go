package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	answer = `
	Here is an information about the authors of this project:

	1. Denis Kondratov. BAI at Yakov & Partners Consulting Company. Student of MIPT DIHT
	2. Aleksei Ovchenkov. Analyst at Yandex.Search Company. Student of MIPT DIHT
	3. Steve Karpov. Backend Developer at Yandex.Taxi Company. Student of MIPT DIHT
	`
)

func SendMessage(ctx context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	fmt.Println("SendMessage Request: ", req)

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
