package execution

import (
	"adk/a2a/server"
	"adk/providers"
	ctx "context"
	"fmt"

	"github.com/google/uuid"
)

func FailAndSaveTask(context ctx.Context, task *server.Task) error {
	fmt.Printf("Task failed: %v\n", task.Id)

	task.Status = server.TaskState_TASK_STATE_FAILED
	// TODO: Save task to storage
	return nil
}

func CreateNewDetachedTask(context ctx.Context, message *server.Message, provider providers.Provider, database, collection string) (*server.SendMessageResponse, error) {
	// 0. Generating a unique task ID
	taskId := uuid.New().String()
	fmt.Println("Task initialization started: ", taskId)

	// 1. Initializing task from the request
	task := InitTask(context, message, taskId)

	// 2. Saving initial state of the task to storage
	err := SaveInitialState(context, task, database, collection)
	if err != nil {
		FailAndSaveTask(context, task)
		return nil, err
	}
	// Create a detached context for the goroutine that won't be canceled
	// when the parent context is canceled
	detachedCtx := ctx.Background()

	go func() {

		// 3. Fetching answers from the task
		history, err := FetchHistory(detachedCtx, task)
		if err != nil {
			FailAndSaveTask(detachedCtx, task)
			return
		}
		// 4. Iterating over answers, getting some agents responses
		task, err := IterateOverAnswers(detachedCtx, provider, task, history)
		if err != nil {
			FailAndSaveTask(detachedCtx, task)
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
