# A2A API

gRPC API для agent-to-agent коммуникации.

## Структура

- `adk/a2a/api/` - proto файлы с описанием API
- `adk/a2a/server/` - сгенерированные .pb.go файлы

## Генерация кода

После изменения `a2a.proto` выполните:

```bash
cd adk
export PATH=$PATH:$(go env GOPATH)/bin
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --go_opt=Ma2a/api/a2a.proto=adk/a2a/server \
       --go-grpc_opt=Ma2a/api/a2a.proto=adk/a2a/server \
       a2a/api/a2a.proto

# Переместить сгенерированные файлы в server
mv a2a/api/*.pb.go a2a/server/
```
