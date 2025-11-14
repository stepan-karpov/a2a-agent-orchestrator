# ADK - Agent Development Kit

ADK (Agent Development Kit) is a Go library for building agent-to-agent orchestrator services. It provides a complete framework for creating gRPC-based services that handle agent communication, task management, LLM provider integration, and persistent storage.

## Overview

ADK simplifies the creation of orchestrator services by providing:

- **gRPC Server Wrapper**: Ready-to-use gRPC server with A2A (Agent-to-Agent) protocol
- **LLM Provider Abstraction**: Unified interface for multiple LLM providers (Eliza, OpenRouter/DeepSeek, etc.)
- **Specialized Agent Support**: Built-in support for registering and calling specialized agents via gRPC
- **Task Management**: Automatic task lifecycle management with MongoDB persistence
- **History Management**: Built-in conversation history tracking and context management
- **Execution Engine**: Detached task execution with automatic state management

## Purpose

ADK is designed to help developers quickly build orchestrator services that:

1. Receive messages from agents or clients via gRPC
2. Process messages using LLM providers
3. Route queries to specialized agents when needed
4. Manage task state and conversation history
5. Persist data to MongoDB
6. Return responses asynchronously

Instead of implementing gRPC servers, task management, provider integrations, and agent coordination from scratch, developers can use ADK to focus on business logic.

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
├── agents/           # Agent management and communication
│   ├── agent.go      # Agent type definition
│   └── call.go       # Agent gRPC client implementation
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
   - Adds system prompt and agent descriptions (if agents are registered)
   - Calls LLM provider with history + new message
   - If LLM requests an agent call, routes query to specialized agent via gRPC
   - Incorporates agent response into conversation history
   - Saves final state with artifacts
4. **GetTask** can be polled to check task status

### Agent Integration

ADK supports registering specialized agents that can be called by the LLM:

1. **Register Agents**: Use `server.RegisterNewAgent()` to register agents with name, description, and gRPC URL
2. **Agent Tools**: Registered agents are automatically exposed to LLM providers as tools
3. **Automatic Routing**: When LLM decides to call an agent, ADK handles the gRPC communication
4. **Context Preservation**: Agent responses are added to conversation history for context

Example:

```go
server.RegisterNewAgent(agents.Agent{
    Name:        "weather-agent",
    Description: "Weather agent. You can use this agent to get information about weather.",
    Url:         "localhost:50056",
})
```

The agent URL is automatically normalized:
- `localhost` or `127.0.0.1` → `passthrough:///localhost:port`
- Other addresses → `dns:///host:port`

### Execution Settings

Pass `ExecutionSettings` to customize behavior:

```go
type ExecutionSettings struct {
    Prompt       string `json:"prompt"`        // System prompt
    HistoryLimit int    `json:"history_limit"` // Max history messages
}
```

### LLM Provider Tool Support

LLM providers (Eliza, OpenRouter) automatically support agent tools:

- **Tool Generation**: Each registered agent becomes a tool with:
  - `name`: Agent name (e.g., "weather-agent")
  - `description`: Agent description (shown to LLM)
  - `parameters`: Single parameter "message" with description "Message to send to remote A2A agent"

- **Tool Calls**: When LLM decides to use an agent:
  - Provider extracts agent name and message from tool call
  - ADK routes the query to the appropriate agent via gRPC
  - Agent response is incorporated into conversation history
  - LLM can continue processing with agent context

## Examples

See `services/orchestrator/` for a complete example of an orchestrator service built with ADK.
