# ADK - Agent Development Kit

ADK (Agent Development Kit) is a Go library for building agent-to-agent orchestrator services. It provides a complete framework for creating gRPC-based services that handle agent communication, task management, LLM provider integration, and persistent storage.

## Overview

ADK simplifies the creation of orchestrator services by providing:

- **gRPC Server Wrapper**: Ready-to-use gRPC server with A2A (Agent-to-Agent) protocol
- **LLM Provider Abstraction**: Unified interface for multiple LLM providers (Eliza, OpenRouter/DeepSeek, etc.)
- **Task Management**: Automatic task lifecycle management with MongoDB persistence
- **History Management**: Built-in conversation history tracking and context management
- **Execution Engine**: Detached task execution with automatic state management

## Purpose

ADK is designed to help developers quickly build orchestrator services that:

1. Receive messages from agents or clients via gRPC
2. Process messages using LLM providers
3. Manage task state and conversation history
4. Persist data to MongoDB
5. Return responses asynchronously

Instead of implementing gRPC servers, task management, and provider integrations from scratch, developers can use ADK to focus on business logic.

## API Reference

The ADK API is defined using Protocol Buffers. The main service contract is in:

- **Proto Definition**: `adk/a2a/api/a2a.proto`
- **Generated Code**: `adk/a2a/server/` (auto-generated `.pb.go` files)

### Service Definition

The `A2AService` provides two main RPCs:

```protobuf
service A2AService {
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
  rpc GetTask(GetTaskRequest) returns (Task);
}
```

### Key Messages

- **`Message`**: Represents a message with context_id, role, content, and metadata
- **`Task`**: Represents a task with id, context_id, status, artifacts, and metadata
- **`Artifact`**: Represents the output/result of a task

See `adk/a2a/api/a2a.proto` for the complete API definition.

## Dependencies

ADK relies on the following technologies and libraries:

### Core Dependencies

- **gRPC** (`google.golang.org/grpc`): For RPC communication
- **Protocol Buffers** (`google.golang.org/protobuf`): For message serialization
- **MongoDB Driver** (`go.mongodb.org/mongo-driver`): For persistent storage
- **UUID** (`github.com/google/uuid`): For task ID generation

### Infrastructure Requirements

- **MongoDB**: For task and conversation history storage
- **LLM Provider APIs**:
  - Eliza API (`https://api.eliza.yandex.net`)
  - OpenRouter API (`https://openrouter.ai/api/v1`) for DeepSeek and other models

### Configuration

ADK uses a `vault.json` file in the project root for secrets:

```json
{
  "eliza-token": "your-eliza-token",
  "openrouter-api-key": "your-openrouter-key",
  "mongodb-uri": "mongodb://user:password@host:port",
  "telegram-bot-token": "your-telegram-token"
}
```

## Architecture

### Components

```
adk/
├── a2a/              # gRPC API definitions and generated code
│   ├── api/          # Protocol Buffer definitions
│   └── server/        # Generated Go code
├── execution/         # Task execution engine
│   ├── execute.go    # Main execution logic
│   ├── 1-init-task.go
│   ├── 2-save-initial-state.go
│   ├── 3-fetch-history.go
│   ├── 4-iterate-over-answers.go
│   └── 5-save-final-state.go
├── providers/         # LLM provider implementations
│   ├── base_provider.go
│   ├── eliza/
│   └── openrouter/
├── storages/         # Storage backends
│   └── mongo/        # MongoDB implementation
├── secrets/          # Secret management
└── server.go         # Main server wrapper
```

### Execution Flow

1. **SendMessage** handler receives a message
2. **CreateNewDetachedTask** creates a task and saves initial state
3. **Detached goroutine** executes:
   - Fetches conversation history
   - Calls LLM provider with history + new message
   - Saves final state with artifacts
4. **GetTask** can be polled to check task status

### Execution Settings

Pass `ExecutionSettings` to customize behavior:

```go
type ExecutionSettings struct {
    Prompt       string `json:"prompt"`        // System prompt
    HistoryLimit int    `json:"history_limit"` // Max history messages
}
```

## Examples

See `services/orchestrator/` for a complete example of an orchestrator service built with ADK.
