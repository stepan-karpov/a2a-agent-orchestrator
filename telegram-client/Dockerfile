FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy entire project directory
COPY . .

# Navigate to telegram-client directory
WORKDIR /app/telegram-client

# Build
RUN go build -o telegram-client .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/telegram-client/telegram-client .

CMD ["./telegram-client"]

