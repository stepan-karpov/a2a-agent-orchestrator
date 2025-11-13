# A2A API

gRPC API for agent-to-agent communication.

## Structure

- `adk/a2a/api/` - proto files with API definitions
- `adk/a2a/server/` - generated .pb.go files

## Code Generation

After modifying `a2a.proto`, run:

```bash
cd adk
export PATH=$PATH:$(go env GOPATH)/bin
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --go_opt=Ma2a/api/a2a.proto=adk/a2a/server \
       --go-grpc_opt=Ma2a/api/a2a.proto=adk/a2a/server \
       a2a/api/a2a.proto

# Move generated files to server
mv a2a/api/*.pb.go a2a/server/
```
