package adk

import (
	"context"
	"fmt"
	"log"
	"net"

	"adk/a2a/server"
	"adk/providers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type SendMessageHandler func(ctx context.Context, req *server.SendMessageRequest, provider providers.Provider) (*server.SendMessageResponse, error)
type GetTaskHandler func(ctx context.Context, req *server.GetTaskRequest) (*server.Task, error)

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
	a2aServer  server.A2AServiceServer
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

	a2aServer := &serverWrapper{
		UnimplementedA2AServiceServer: server.UnimplementedA2AServiceServer{},
		sendMessageHandler:            config.SendMessageHandler,
		getTaskHandler:                config.GetTaskHandler,
		provider:                      config.Provider,
	}

	return &Server{
		config:    config,
		a2aServer: a2aServer,
		provider:  config.Provider,
	}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.config.Port, err)
	}

	s.grpcServer = grpc.NewServer()
	server.RegisterA2AServiceServer(s.grpcServer, s.a2aServer)

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
