package execution

import (
	"adk/a2a/server"
	"adk/storages/mongo"
	"context"
)

func SaveInitialState(context context.Context, task *server.Task, database, collection string) error {
	err := mongo.SaveTask(context, task, database, collection)
	if err != nil {
		return err
	}
	return nil
}
