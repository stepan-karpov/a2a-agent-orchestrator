FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy entire project directory
COPY . .

# Navigate to agent directory
WORKDIR /app/services/agents/crypto-agent

# Build (disable workspace to avoid go.work conflicts)
ENV GOWORK=off
RUN go build -o crypto-agent .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/agents/crypto-agent/crypto-agent .

CMD ["./crypto-agent"]

