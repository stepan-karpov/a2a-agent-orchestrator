FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy entire project directory
COPY . .

# Navigate to agent directory
WORKDIR /app/services/agents/weather-agent

# Build
RUN go build -o weather-agent .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/agents/weather-agent/weather-agent .

CMD ["./weather-agent"]

