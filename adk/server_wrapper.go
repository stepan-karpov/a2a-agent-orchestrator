package adk

import (
	"adk/a2a/server"
	"adk/providers"
	"context"
)

type serverWrapper struct {
	server.UnimplementedA2AServiceServer
	sendMessageHandler SendMessageHandler
	getTaskHandler     GetTaskHandler
	provider           providers.Provider
}

func (w *serverWrapper) SendMessage(ctx context.Context, req *server.SendMessageRequest) (*server.SendMessageResponse, error) {
	return w.sendMessageHandler(ctx, req, w.provider)
}

func (w *serverWrapper) GetTask(ctx context.Context, req *server.GetTaskRequest) (*server.Task, error) {
	return w.getTaskHandler(ctx, req)
}
