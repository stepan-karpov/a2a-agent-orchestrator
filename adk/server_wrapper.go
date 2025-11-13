package adk

import (
	a2aServerProto "adk/a2a/server"
	"context"
)

type serverWrapper struct {
	a2aServerProto.UnimplementedA2AServiceServer
	sendMessageHandler SendMessageHandler
	getTaskHandler     GetTaskHandler
	server             *Server
}

func (w *serverWrapper) SendMessage(context context.Context, req *a2aServerProto.SendMessageRequest) (*a2aServerProto.SendMessageResponse, error) {
	return w.sendMessageHandler(context, req, w.server)
}

func (w *serverWrapper) GetTask(context context.Context, req *a2aServerProto.GetTaskRequest) (*a2aServerProto.Task, error) {
	return w.getTaskHandler(context, req, w.server)
}
