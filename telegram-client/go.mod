module telegram-client

go 1.25.3

require (
	adk v0.0.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace adk => ../adk
