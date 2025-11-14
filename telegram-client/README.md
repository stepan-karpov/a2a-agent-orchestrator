# Telegram Client

Telegram bot client that connects users to the A2A orchestrator service. This service receives messages from Telegram users, forwards them to the orchestrator via gRPC, polls for task completion, and sends responses back to users.

## Overview

The Telegram client acts as a bridge between Telegram users and the orchestrator service:

- Receives messages from Telegram users
- Provides an interactive start menu with quick access buttons for common topics
- Sends messages to the orchestrator via gRPC (`SendMessage` RPC)
- Polls task status until completion or failure
- Returns the response to the user in Telegram

## Features

- **Interactive Start Menu**: Users receive a welcome message with inline keyboard buttons when they send `/start`
- **Quick Access Buttons**: Pre-configured buttons for:
  - About authors of this project
  - About crypto
  - About weather
  - About news
  - About football
- **Specialized Agents**: The bot can route queries to specialized agents through the orchestrator
- **HTML Formatting**: Supports Telegram HTML formatting in responses

## Project Structure

```
telegram-client/
├── main.go              # Main bot logic and message handling
├── vault.go             # Secret management (reads vault.json)
├── Dockerfile           # Container build configuration
├── docker-compose.yml   # Docker Compose configuration
├── go.mod               # Go module dependencies
└── README.md            # This file
```

## How It Works

1. **Start Command**: When user sends `/start`, bot displays welcome message with inline keyboard buttons
2. **Button Interaction**: User can click buttons or type messages directly
3. **Message Reception**: Bot listens for Telegram messages and callback queries via long polling
4. **gRPC Request**: Sends message to orchestrator's `SendMessage` endpoint
5. **Task Polling**: Polls `GetTask` endpoint up to 50 times with increasing delays
6. **Response**: Sends the task result (artifact content) back to the user

## Prerequisites

To run the Telegram client, you need:

1. **Orchestrator service** - Must be running and accessible

   - Default gRPC endpoint: `localhost:50051` (for local development)
   - For Docker: configure `GRPC_HOST` environment variable

2. **Telegram Bot Token** - Required in `vault.json`

   - Get a token from [@BotFather](https://t.me/botfather) on Telegram
   - Add it to `vault.json` as `telegram-bot-token`

3. **vault.json** - Configuration file with secrets
   - Must be in the project root directory
   - Must contain `telegram-bot-token` field

## Configuration

### Environment Variables

- `GRPC_HOST` - gRPC endpoint of the orchestrator service (default: `localhost:50051`)
- `TELEGRAM_BOT_TOKEN` - Telegram bot token (optional, can be read from vault.json)

### vault.json

The client reads the Telegram bot token from `vault.json` in the project root:

```json
{
  "telegram-bot-token": "YOUR_BOT_TOKEN_HERE"
}
```

## Local Development

### Run Locally

```bash
cd telegram-client
go run .
```

Make sure:

- Orchestrator is running on `localhost:50051`
- `vault.json` exists in the project root with `telegram-bot-token`

## Deployment

### Docker Compose

1. Update `docker-compose.yml` with the correct `GRPC_HOST`:

   ```yaml
   environment:
     - GRPC_HOST=orchestrator:50051  # or your orchestrator host
   ```

2. Build and run:

   ```bash
   cd telegram-client
   podman-compose up -d --build
   # or
   docker-compose up -d --build
   ```

3. Check logs:
   ```bash
   podman logs telegram-client
   # or
   docker logs telegram-client
   ```

### Docker (Standalone)

1. Build the image:

   ```bash
   docker build -f telegram-client/Dockerfile -t telegram-client:latest ..
   ```

2. Run the container:
   ```bash
   docker run -d \
     --name telegram-client \
     -e GRPC_HOST=orchestrator:50051 \
     -v $(pwd)/vault.json:/app/vault.json:ro \
     telegram-client:latest
   ```

### Network Configuration

For Docker deployments, ensure the Telegram client can reach the orchestrator:

- **Same Docker network**: Use the same network or connect containers to the same network
- **Host network**: Use `network_mode: host` in docker-compose.yml (works for Podman)
- **Direct IP**: Set `GRPC_HOST` to the orchestrator's IP address

Example for Podman with host network:

```yaml
network_mode: host
environment:
  - GRPC_HOST=127.0.0.1:50051
```

