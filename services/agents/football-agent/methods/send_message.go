package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"

	"github.com/google/uuid"
)

const answer = `Here is the information about the football matches:
1. Manchester United vs. Liverpool
2. Chelsea vs. Arsenal
3. Manchester City vs. Tottenham
4. Barcelona vs. Real Madrid
5. Bayern Munich vs. Borussia Dortmund
6. Juventus vs. Inter Milan
7. Paris Saint-Germain vs. Lyon
8. Atletico Madrid vs. Real Sociedad`

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
