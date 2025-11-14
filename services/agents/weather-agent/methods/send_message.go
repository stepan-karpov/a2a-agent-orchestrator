package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"
	"fmt"

	"github.com/google/uuid"
)

const answer = `Weather forecast for today:
Temperature: 22°C
Condition: Partly cloudy
Humidity: 65%
Wind: 15 km/h
Precipitation: 10% chance of rain`

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
