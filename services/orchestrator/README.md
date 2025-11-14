# Orchestrator Service

The Orchestrator service is an agent-to-agent communication service built using the ADK (Agent Development Kit) library. It provides a gRPC API for processing messages, managing tasks, and coordinating agent interactions using LLM providers.

## Overview

This service acts as the central orchestrator for agent-to-agent communication:

- Receives messages via gRPC (`SendMessage` RPC)
- Processes messages using LLM providers (Eliza, OpenRouter, etc.)
- Manages task lifecycle and conversation history
- Coordinates with specialized agents (authors, crypto, weather, news, football)
- Persists data to MongoDB
- Returns task results asynchronously

## How It Works

1. **Message Reception**: Client sends message via `SendMessage` RPC
2. **Task Creation**: Handler creates a detached task using `server.CreateNewDetachedTask()`
3. **Background Processing**:
   - Fetches conversation history from MongoDB
   - Adds system prompt (from `orchestrator_setting.json`)
   - Limits history to `history_limit` messages
   - Calls LLM provider with full context
   - Saves result as artifact
4. **Task Completion**: Task status changes to `TASK_STATE_COMPLETED`
5. **Polling**: Client polls `GetTask` to retrieve results

## Architecture

The orchestrator is built on top of ADK and implements:

- **gRPC Server**: Exposes A2A API on port `:50051`
- **LLM Provider Integration**: Uses ADK's provider abstraction (default: Eliza)
- **Task Management**: Automatic task creation, execution, and state management
- **History Management**: Conversation context tracking with configurable limits
- **MongoDB Storage**: Persistent storage for tasks and conversation history

## Prerequisites

Before running the orchestrator, ensure you have:

1. **MongoDB** - Running and accessible

   - Connection string configured in `vault.json` as `mongodb-uri`
   - Default database: `a2a`
   - Default collection: `orchestrator`

2. **LLM Provider Credentials** - In `vault.json`:

   - `eliza-token` - For Eliza provider (default)
   - `openrouter-api-key` - For OpenRouter/DeepSeek provider

3. **vault.json** - In project root with required secrets:

   ```json
   {
     "eliza-token": "your-eliza-token",
     "openrouter-api-key": "your-openrouter-key",
     "mongodb-uri": "mongodb://user:password@host:port"
   }
   ```

4. **orchestrator_setting.json** - Service configuration:
   ```json
   {
     "prompt": "You are a Telegram bot assistant...",
     "history_limit": 10,
     "agents": [
       {
         "name": "authors-info-agent",
         "description": "Authors info agent...",
         "url": "localhost:50052"
       },
       {
         "name": "crypto-agent",
         "description": "Crypto agent...",
         "url": "localhost:50053"
       }
       // ... more agents
     ]
   }
   ```
   
   The `agents` array defines specialized agents that the orchestrator can route queries to. Each agent has:
   - `name`: Agent identifier used in tool calls
   - `description`: Description shown to LLM for agent selection
   - `url`: gRPC endpoint where the agent service is running

## Configuration

### orchestrator_setting.json

Located in `services/orchestrator/orchestrator_setting.json`:

- **`prompt`**: System prompt sent to LLM as the first message (role: "system")
- **`history_limit`**: Maximum number of conversation messages to include in context (excluding system prompt)
- **`agents`**: Array of specialized agents that can be called by the LLM. Each agent includes:
  - `name`: Unique identifier for the agent
  - `description`: Description of what the agent does (shown to LLM)
  - `url`: gRPC endpoint (e.g., `localhost:50052`)
  
  The orchestrator automatically registers these agents and makes them available to the LLM as tools. When the LLM determines that an agent should be called, the orchestrator routes the query to the appropriate agent service.

### Server Configuration

Configured in `main.go`:

- **Port**: `:50051` (gRPC server port)
- **Provider**: `"eliza"` (can be changed to `"openrouter"`)
- **Database**: `"a2a"`
- **Collection**: `"orchestrator"`

## Running the Service

### Local Development

```bash
cd services/orchestrator
go run main.go
```

Or from project root:

```bash
go run services/orchestrator/main.go
```

### Using Go Workspace

```bash
# From project root
go run ./services/orchestrator
```

## API Usage

### SendMessage

Sends a message to the orchestrator and creates a task:

```bash
grpcurl -plaintext -d '{
  "request": {
    "context_id": "my-context-id-1",
    "role": "ROLE_USER",
    "content": "Hello, how are you?"
  }
}' localhost:50051 a2a.A2AService/SendMessage
```

Response:

```json
{
  "task": {
    "id": "task-uuid",
    "context_id": "my-context-id-1",
    "status": "TASK_STATE_WORKING"
  }
}
```

### GetTask

Retrieves task status and results:

```bash
grpcurl -plaintext -d '{
  "name": "task-uuid-here"
}' localhost:50051 a2a.A2AService/GetTask
```

Response:

```json
{
  "id": "task-uuid",
  "context_id": "my-context-id-1",
  "status": "TASK_STATE_COMPLETED",
  "artifacts": [
    {
      "type": "text",
      "content": "LLM response here"
    }
  ]
}
```

## Project Structure

```
services/orchestrator/
├── main.go                    # Service entry point
├── methods/                   # gRPC handler implementations
│   ├── send_message.go       # SendMessage handler
│   └── get_task.go           # GetTask handler
├── orchestrator_setting.json  # Service configuration
├── go.mod                     # Go module dependencies
└── README.md                  # This file
```

## Dependencies

- **ADK Library**: Core orchestrator framework
- **MongoDB Driver**: For persistent storage
- **gRPC**: For API communication
- **Protocol Buffers**: For message serialization

See `go.mod` for complete dependency list.
