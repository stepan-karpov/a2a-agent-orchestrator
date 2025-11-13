package adk

import (
	"context"
	"fmt"
	"log"
	"net"

	a2aServerProto "adk/a2a/server"
	"adk/execution"
	"adk/providers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/structpb"
)

type SendMessageHandler func(context context.Context, req *a2aServerProto.SendMessageRequest, server *Server) (*a2aServerProto.SendMessageResponse, error)
type GetTaskHandler func(context context.Context, req *a2aServerProto.GetTaskRequest) (*a2aServerProto.Task, error)

type ServerConfig struct {
	Port               string
	Provider           providers.Provider
	SendMessageHandler SendMessageHandler
	GetTaskHandler     GetTaskHandler
}

// Server - adk wrapper for gRPC server
type Server struct {
	config     ServerConfig
	grpcServer *grpc.Server
	a2aServer  a2aServerProto.A2AServiceServer
	provider   providers.Provider
}

// NewServer creates a new server instance with configuration
func NewServer(config ServerConfig) (*Server, error) {
	if config.Provider == nil {
		return nil, fmt.Errorf("provider is required")
	}
	if config.SendMessageHandler == nil {
		return nil, fmt.Errorf("SendMessageHandler is required")
	}
	if config.GetTaskHandler == nil {
		return nil, fmt.Errorf("GetTaskHandler is required")
	}

	server := &Server{
		config:   config,
		provider: config.Provider,
	}

	a2aServer := &serverWrapper{
		UnimplementedA2AServiceServer: a2aServerProto.UnimplementedA2AServiceServer{},
		sendMessageHandler:            config.SendMessageHandler,
		getTaskHandler:                config.GetTaskHandler,
		server:                        server,
	}

	server.a2aServer = a2aServer

	return server, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.config.Port, err)
	}

	s.grpcServer = grpc.NewServer()
	a2aServerProto.RegisterA2AServiceServer(s.grpcServer, s.a2aServer)

	reflection.Register(s.grpcServer)

	log.Printf("gRPC server listening on %s", s.config.Port)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

type UserMessageInfo struct {
	MessageId string
	ContextId string
	TaskId    string
	Role      string
	Content   string
	Metadata  *structpb.Struct
}

// Creates a new detached task for ADK execution
func (s *Server) CreateNewDetachedTask(context context.Context, message *a2aServerProto.Message) (*a2aServerProto.Task, error) {
	response, err := execution.CreateNewDetachedTask(context, message, s.provider)
	if err != nil {
		return nil, err
	}
	return response.Task, nil
}
