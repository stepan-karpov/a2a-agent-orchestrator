FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy entire project directory
COPY . .

# Navigate to agent directory
WORKDIR /app/services/rag

# Build (disable workspace to avoid go.work conflicts)
ENV GOWORK=off
# Download dependencies and build (skip verify for local adk module)
# Build with explicit output name
RUN go mod download && go build -o rag . && ls -la

FROM alpine:latest

WORKDIR /app

# Copy binary and required files
# Check which binary was created (rag or rag-platform)
COPY --from=builder /app/services/rag/rag* /app/
COPY --from=builder /app/services/rag/settings.json /app/settings.json
COPY --from=builder /app/services/rag/database.xlsx /app/database.xlsx

# Make binary executable and create symlink if needed
RUN chmod +x /app/rag* && \
    if [ -f /app/rag-platform ] && [ ! -f /app/rag ]; then \
        mv /app/rag-platform /app/rag; \
    fi

CMD ["/app/rag"]

