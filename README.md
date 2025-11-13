# A2A Agent Orchestrator

A complete agent-to-agent communication platform built with Go, providing a framework for orchestrating LLM-powered agents, managing conversations, and handling asynchronous task processing.

## Overview

This project consists of three main components:

1. **ADK (Agent Development Kit)** - A Go library for building orchestrator services
2. **Orchestrator Service** - A reference implementation using ADK
3. **Telegram Client** - A Telegram bot that connects users to the orchestrator

The platform enables developers to quickly build services that process messages using LLM providers, manage conversation history, and handle asynchronous task execution with persistent storage.

## Architecture

```
┌─────────────────┐
│ Telegram Client │  (User Interface)
└────────┬────────┘
         │ gRPC
         ▼
┌─────────────────┐
│  Orchestrator   │  (Built with ADK)
│     Service     │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌────────┐ ┌──────────┐
│ MongoDB│ │ LLM APIs │
│(Storage)│ │(Eliza,   │
│        │ │OpenRouter)│
└────────┘ └──────────┘
```

### Component Overview

- **ADK**: Core library providing gRPC server wrapper, LLM provider abstraction, task management, and MongoDB integration
- **Orchestrator**: Production-ready service that processes messages, manages tasks, and maintains conversation context
- **Telegram Client**: User-facing bot that bridges Telegram users with the orchestrator service

## Project Structure

```
a2a-agent-orchestrator/
├── adk/                    # Agent Development Kit library
│   ├── a2a/               # gRPC API definitions
│   ├── execution/         # Task execution engine
│   ├── providers/         # LLM provider implementations
│   ├── storages/          # Storage backends (MongoDB)
│   └── README.md          # ADK documentation
├── services/
│   ├── orchestrator/      # Orchestrator service implementation
│   │   └── README.md      # Orchestrator documentation
│   └── simple-agent/      # Additional agent service
├── telegram-client/       # Telegram bot client
│   └── README.md          # Telegram client documentation
├── vault.json             # Secrets configuration (not in repo)
└── go.work                # Go workspace configuration
```

## Quick Start

### Prerequisites

- Go 1.25.3 or later
- MongoDB (running and accessible)
- LLM provider credentials (Eliza token or OpenRouter API key)
- Telegram bot token (for Telegram client)

### Setup

1. **Configure secrets** - Create `vault.json` in project root:

   ```json
   {
     "eliza-token": "your-token",
     "openrouter-api-key": "your-key",
     "mongodb-uri": "mongodb://user:password@host:port",
     "telegram-bot-token": "your-bot-token"
   }
   ```

2. **Start MongoDB** - Ensure MongoDB is running and accessible

3. **Run orchestrator**:

   ```bash
   go run ./services/orchestrator
   ```

4. **Run Telegram client** (optional):
   ```bash
   go run ./telegram-client
   ```

## Documentation

- **[ADK Library](adk/README.md)** - Framework documentation, API reference, and integration guide
- **[Orchestrator Service](services/orchestrator/README.md)** - Service setup, configuration, and API usage
- **[Telegram Client](telegram-client/README.md)** - Bot setup, deployment, and troubleshooting
- **[A2A API](adk/a2a/README.md)** - Protocol Buffer definitions and code generation

## Key Features

- **Asynchronous Processing**: Tasks are processed in background goroutines, allowing non-blocking request handling
- **Conversation History**: Automatic tracking and retrieval of conversation context per `context_id`
- **Multiple LLM Providers**: Support for Eliza, OpenRouter/DeepSeek, and extensible provider interface
- **Persistent Storage**: MongoDB integration for task and conversation history storage
- **gRPC API**: Standardized A2A (Agent-to-Agent) protocol for service communication
- **Configurable Prompts**: System prompts and history limits configurable per service

## Technology Stack

- **Language**: Go 1.25.3
- **Communication**: gRPC with Protocol Buffers
- **Storage**: MongoDB
- **LLM Providers**: Eliza API, OpenRouter API
- **Deployment**: Docker/Podman support

## Development

This project uses Go workspaces. All modules are defined in `go.work`:

- `./adk` - Core library module
- `./services/orchestrator` - Orchestrator service module
- `./services/simple-agent` - Simple agent module
- `./telegram-client` - Telegram client module

### Building

```bash
# Build all modules
go build ./...

# Build specific service
go build ./services/orchestrator
go build ./telegram-client
```

### Testing

Use `grpcurl` to test the orchestrator API:

```bash
# Send message
grpcurl -plaintext -d '{
  "request": {
    "context_id": "test-context",
    "role": "ROLE_USER",
    "content": "Hello"
  }
}' localhost:50051 a2a.A2AService/SendMessage
```
