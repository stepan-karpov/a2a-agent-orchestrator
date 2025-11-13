package execution

import (
	"adk/a2a/server"
	"adk/providers"
	"context"
)

func FetchHistory(context context.Context, task *server.Task) ([]providers.Message, error) {
	return []providers.Message{}, nil
}
