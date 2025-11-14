package methods

import (
	"adk"
	a2aServerProto "adk/a2a/server"
	"context"
	"fmt"

	"github.com/google/uuid"
)

const answer = `Current cryptocurrency prices:
Bitcoin (BTC): $43,250.00 (+2.5%)
Ethereum (ETH): $2,680.00 (+1.8%)
Binance Coin (BNB): $315.00 (+0.9%)
Solana (SOL): $98.50 (+3.2%)
Cardano (ADA): $0.52 (+1.1%)

Market cap: $1.65T
24h volume: $85.2B`

func SendMessage(ctx context.Context, req *a2aServerProto.SendMessageRequest, server *adk.Server) (*a2aServerProto.SendMessageResponse, error) {
	fmt.Println("SendMessage Request: ", req)

	task := &a2aServerProto.Task{
		Id:        uuid.New().String(),
		ContextId: req.Request.ContextId,
		Status:    a2aServerProto.TaskState_TASK_STATE_COMPLETED,
		Artifacts: []*a2aServerProto.Artifact{
			{
				Type:    "text",
				Content: answer,
			},
		},
	}

	return &a2aServerProto.SendMessageResponse{Task: task}, nil
}
