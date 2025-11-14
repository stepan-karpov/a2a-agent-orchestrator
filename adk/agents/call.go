package agents

import (
	"context"
	"fmt"
	"strings"
	"time"

	a2aProto "adk/a2a/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CallAgent(ctx context.Context, agent *Agent, message string, contextId string) (string, error) {
	// Ensure URL has a resolver scheme for grpc.NewClient
	// Use passthrough for localhost, dns for other addresses
	url := agent.Url
	if !strings.Contains(url, "://") {
		if strings.HasPrefix(url, "localhost") || strings.HasPrefix(url, "127.0.0.1") {
			url = "passthrough:///" + url
		} else {
			url = "dns:///" + url
		}
	}

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	if req.Request == nil {
		return "", fmt.Errorf("failed to create message request: request is nil")
	}

	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to send message to agent: %w", err)
	}

	if resp == nil || resp.Task == nil {
		return "", fmt.Errorf("invalid response from agent: task is nil")
	}

	taskId := resp.Task.Id

	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond * time.Duration(i))

		getReq := &a2aProto.GetTaskRequest{Name: taskId}
		task, err := client.GetTask(ctx, getReq)
		if err != nil {
			return "", fmt.Errorf("failed to get task status: %w", err)
		}

		if task == nil {
			return "", fmt.Errorf("invalid task response: task is nil")
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
