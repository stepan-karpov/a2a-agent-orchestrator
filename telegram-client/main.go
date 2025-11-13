package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	a2aProto "adk/a2a/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getGRPCHost() string {
	if host := os.Getenv("GRPC_HOST"); host != "" {
		return host
	}
	return "localhost:50051"
}

func ReplyForAMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Connect to gRPC service for each request
	conn, err := grpc.NewClient(getGRPCHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при подключении к сервису-оркестратору")
		bot.Send(msg)
		return
	}
	defer conn.Close()

	grpcClient := a2aProto.NewA2AServiceClient(conn)

	// Send message to gRPC service
	ctx := context.Background()
	req := &a2aProto.SendMessageRequest{
		Request: &a2aProto.Message{
			ContextId: fmt.Sprintf("tg-%d", message.Chat.ID),
			Role:      a2aProto.Role_ROLE_USER,
			Content:   message.Text,
		},
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Обрабатываю ваше сообщение...")
	bot.Send(msg)

	resp, err := grpcClient.SendMessage(ctx, req)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error: "+err.Error())
		bot.Send(msg)
		return
	}

	taskId := resp.Task.Id

	// Poll task status up to 10 times
	var task *a2aProto.Task
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond * time.Duration(i))

		getReq := &a2aProto.GetTaskRequest{Name: taskId}
		task, err = grpcClient.GetTask(ctx, getReq)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при получении статуса задачи: "+err.Error())
			bot.Send(msg)
		}
		if task.Status == a2aProto.TaskState_TASK_STATE_COMPLETED || task.Status == a2aProto.TaskState_TASK_STATE_FAILED {
			break
		}
	}

	// Send response
	responseText := "Processing..."
	if task == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка - задача потерялась в оркестраторе :(")
		bot.Send(msg)
		return
	}

	if task.Status == a2aProto.TaskState_TASK_STATE_COMPLETED && len(task.Artifacts) > 0 {
		responseText = task.Artifacts[0].Content
	} else if task.Status == a2aProto.TaskState_TASK_STATE_FAILED {
		responseText = "Задача упала :("
	}
	msg = tgbotapi.NewMessage(message.Chat.ID, responseText)
	bot.Send(msg)
}

func Listen(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		ReplyForAMessage(bot, update.Message)
	}
}

func main() {
	token := GetTelegramToken()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	Listen(bot)
}
