package execution

import (
	"adk/a2a/server"
	"adk/providers"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func FailAndSaveTask(context context.Context, task *server.Task) error {
	fmt.Errorf("Task failed: %v", task.Id)

	task.Status = server.TaskState_TASK_STATE_FAILED
	// TODO: Save task to storage
	return nil
}

func CreateNewDetachedTask(context context.Context, message *server.Message, provider providers.Provider) (*server.SendMessageResponse, error) {
	// 0. Generating a unique task ID
	taskId := uuid.New().String()

	// 1. Initializing task from the request
	task := InitTask(context, message, taskId)

	// 2. Saving initial state of the task to storage
	err := SaveInitialState(context, task)
	if err != nil {
		FailAndSaveTask(context, task)
		return nil, err
	}

	go func() {
		// 3. Fetching answers from the task
		history, err := FetchHistory(context, task)
		if err != nil {
			FailAndSaveTask(context, task)
			return
		}
		// 4. Iterating over answers, getting some agents responses
		task, err := IterateOverAnswers(context, provider, task, history)
		if err != nil {
			FailAndSaveTask(context, task)
			return
		}
		fmt.Println("Task answer received: ", task.Id)
		// 5. Saving final state of the task to storage
		err = SaveFinalState(context, task)
		if err != nil {
			fmt.Printf("Error saving final state of the task: %v", err)
			return
		}
		fmt.Println("Task completed: ", task.Id)
	}()

	fmt.Println("Task created: ", task.Id)
	return &server.SendMessageResponse{Task: task}, nil
}
