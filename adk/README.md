```
grpcurl -plaintext -d '{
  "request": {
    "context_id": "my-context-id-1",
    "role": "ROLE_USER",
    "content": "Hello, how are you?"
  }
}' localhost:50051 a2a.A2AService/SendMessage
```

```
grpcurl -plaintext -d '{
  "name": "42f1744b-cc87-402b-9aac-46290b5dec66"
}' localhost:50051 a2a.A2AService/GetTask
```
