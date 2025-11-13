package execution

import (
	"adk/a2a/server"
	"adk/providers"
	"context"
)

func IterateOverAnswers(context context.Context, provider providers.Provider, task *server.Task, history []providers.Message) (*server.Task, error) {
	history = append(history,
		providers.Message{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
	)
	response, err := provider.ChatCompletion(history)

	if err != nil {
		return nil, err
	}

	task.Artifacts = append(task.Artifacts, &server.Artifact{
		Type:    "text",
		Content: response.Content,
	})
	task.History = append(task.History, &server.Message{
		ContextId: task.ContextId,
		Role:      server.Role_ROLE_ASSISTANT,
		Content:   response.Content,
	})
	task.Status = server.TaskState_TASK_STATE_COMPLETED

	return task, nil
}
