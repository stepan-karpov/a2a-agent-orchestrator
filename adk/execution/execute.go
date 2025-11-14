package execution

import (
	"adk/a2a/server"
	"adk/agents"
	"adk/providers"
	"adk/storages/mongo"
	ctx "context"
	"fmt"

	"github.com/google/uuid"
)

type ExecutionSettings struct {
	Prompt       string `json:"prompt"`
	HistoryLimit int    `json:"history_limit"`
}

func FailAndSaveTask(context ctx.Context, task *server.Task, database, collection string) error {
	if task == nil {
		fmt.Printf("Task failed: task is nil\n")
		return nil
	}

	fmt.Printf("Task failed: %v\n", task.Id)

	task.Status = server.TaskState_TASK_STATE_FAILED

	err := mongo.SaveTask(context, task, database, collection)
	if err != nil {
		return fmt.Errorf("failed to save failed task: %w", err)
	}

	return nil
}

func CreateNewDetachedTask(context ctx.Context, message *server.Message, serverSettings *ExecutionSettings, provider providers.Provider, database, collection string, agents []agents.Agent) (*server.SendMessageResponse, error) {
	// 0. Generating a unique task ID
	taskId := uuid.New().String()
	fmt.Printf("Task initialization started: %s, message content='%s' (len=%d)\n", taskId, message.Content, len(message.Content))

	// 1. Initializing task from the request
	task := InitTask(context, message, taskId)

	// 2. Saving initial state of the task to storage
	err := SaveInitialState(context, task, database, collection)
	if err != nil {
		fmt.Printf("Error saving initial state of the task: %v", err)
		FailAndSaveTask(context, task, database, collection)
		return nil, err
	}
	// Create a detached context for the goroutine that won't be canceled
	// when the parent context is canceled
	detachedCtx := ctx.Background()

	go func() {

		// 3. Fetching answers from the task
		history, err := FetchHistory(detachedCtx, task, database, collection)
		if err != nil {
			FailAndSaveTask(detachedCtx, task, database, collection)
			return
		}

		// Take last HistoryLimit messages
		if len(history) > serverSettings.HistoryLimit {
			history = history[len(history)-serverSettings.HistoryLimit:]
		}

		// Add prompt as the first message (system message)
		if serverSettings.Prompt != "" {
			promptContent := serverSettings.Prompt

			// Add agent descriptions to prompt if agents are available
			if len(agents) > 0 {
				promptContent += "\n\nAvailable agents:"
				for _, agent := range agents {
					promptContent += fmt.Sprintf("\n- %s: %s", agent.Name, agent.Description)
				}
			}

			history = append([]providers.Message{
				{
					Role:    "system",
					Content: promptContent,
				},
			}, history...)
		}

		fmt.Println("Number of messages (with prompt) in history: ", len(history))

		// 4. Iterating over answers, getting some agents responses
		task, err := IterateOverAnswers(detachedCtx, provider, task, message, history, agents)
		if err != nil {
			fmt.Printf("Error iterating over answers: %v", err)
			FailAndSaveTask(detachedCtx, task, database, collection)
			return
		}
		fmt.Println("Task answer received: ", task.Id)
		// 5. Saving final state of the task to storage
		err = SaveFinalState(detachedCtx, task, database, collection)
		if err != nil {
			fmt.Printf("Error saving final state of the task: %v", err)
			return
		}
		fmt.Println("Task completed: ", task.Id)
	}()

	fmt.Println("Task created: ", task.Id)
	return &server.SendMessageResponse{Task: task}, nil
}
