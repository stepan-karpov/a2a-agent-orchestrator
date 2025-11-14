FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy entire project directory
COPY . .

# Navigate to agent directory
WORKDIR /app/services/agents/news-agent

# Build (disable workspace to avoid go.work conflicts)
ENV GOWORK=off
RUN go build -o news-agent .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/agents/news-agent/news-agent .

CMD ["./news-agent"]

