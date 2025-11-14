FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy entire project directory
COPY . .

# Navigate to orchestrator directory
WORKDIR /app/services/orchestrator

# Build
RUN go build -o orchestrator .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/orchestrator/orchestrator .
COPY --from=builder /app/services/orchestrator/orchestrator_setting.json .

CMD ["./orchestrator"]

