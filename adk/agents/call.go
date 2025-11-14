package agents

import (
	"context"
	"fmt"
	"time"

	a2aProto "adk/a2a/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CallAgent(ctx context.Context, agent *Agent, message string, contextId string) (string, error) {
	conn, err := grpc.NewClient(agent.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("failed to connect to agent %s at %s: %w", agent.Name, agent.Url, err)
	}
	defer conn.Close()

	client := a2aProto.NewA2AServiceClient(conn)

	req := &a2aProto.SendMessageRequest{
		Request: &a2aProto.Message{
			ContextId: contextId,
			Role:      a2aProto.Role_ROLE_USER,
			Content:   message,
		},
	}

	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to send message to agent: %w", err)
	}

	taskId := resp.Task.Id

	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond * time.Duration(i))

		getReq := &a2aProto.GetTaskRequest{Name: taskId}
		task, err := client.GetTask(ctx, getReq)
		if err != nil {
			return "", fmt.Errorf("failed to get task status: %w", err)
		}

		if task.Status == a2aProto.TaskState_TASK_STATE_COMPLETED {
			if len(task.Artifacts) > 0 {
				return task.Artifacts[0].Content, nil
			}
			return "", fmt.Errorf("task completed but no artifacts found")
		}

		if task.Status == a2aProto.TaskState_TASK_STATE_FAILED {
			return "", fmt.Errorf("agent task failed")
		}
	}

	return "", fmt.Errorf("timeout waiting for agent response")
}
