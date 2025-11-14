package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"
	"fmt"

	"github.com/google/uuid"
)

const answer = `Latest News Headlines:

1. Tech Breakthrough: New AI Model Achieves Human-Level Performance
   Published: 2 hours ago | Source: TechNews

2. Global Markets Rally on Economic Optimism
   Published: 5 hours ago | Source: Financial Times

3. Climate Summit Reaches Historic Agreement
   Published: 8 hours ago | Source: Reuters

4. Space Mission Successfully Launches to Mars
   Published: 12 hours ago | Source: NASA

5. Healthcare Innovation: New Treatment Shows Promise
   Published: 1 day ago | Source: Medical Journal`

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
