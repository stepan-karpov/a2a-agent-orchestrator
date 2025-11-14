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

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(msg)
	return err
}

func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) error {
	text := `Welcome! 👋

I'm an AI assistant that can help you get information on various topics using specialized agents. You can ask me questions or use the buttons below to quickly access information about:

- <b>Authors</b> - Learn about the project authors
- <b>Crypto</b> - Get information about cryptocurrency
- <b>Weather</b> - Check weather information
- <b>News</b> - Read the latest news
- <b>Football</b> - Get football match information

Choose what you'd like to know about:`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("About authors", "ask_authors"),
			tgbotapi.NewInlineKeyboardButtonData("About crypto", "ask_crypto"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("About weather", "ask_weather"),
			tgbotapi.NewInlineKeyboardButtonData("About news", "ask_news"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("About football", "ask_football"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(msg)
	return err
}

func ReplyForAMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Connect to gRPC service for each request
	conn, err := grpc.NewClient(getGRPCHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		SendMessage(bot, message.Chat.ID, "Error connecting to orchestrator service")
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

	SendMessage(bot, message.Chat.ID, "Processing your message...")

	resp, err := grpcClient.SendMessage(ctx, req)
	if err != nil {
		SendMessage(bot, message.Chat.ID, "Error: "+err.Error())
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
			SendMessage(bot, message.Chat.ID, "Error getting task status: "+err.Error())
		}
		if task.Status == a2aProto.TaskState_TASK_STATE_COMPLETED || task.Status == a2aProto.TaskState_TASK_STATE_FAILED {
			break
		}
	}

	// Send response
	responseText := "Processing..."
	if task == nil {
		SendMessage(bot, message.Chat.ID, "Error - task was lost in orchestrator :(")
		return
	}

	if task.Status == a2aProto.TaskState_TASK_STATE_COMPLETED && len(task.Artifacts) > 0 {
		responseText = task.Artifacts[0].Content
	} else if task.Status == a2aProto.TaskState_TASK_STATE_FAILED {
		responseText = "Task failed :("
	}
	SendMessage(bot, message.Chat.ID, responseText)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Acknowledge the callback
	bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Map callback data to questions
	questionMap := map[string]string{
		"ask_authors":  "Tell me about the authors of this project",
		"ask_crypto":   "Tell me about crypto",
		"ask_weather":  "Tell me about weather",
		"ask_news":     "Tell me about news",
		"ask_football": "Tell me about football",
	}

	question, exists := questionMap[callback.Data]
	if !exists {
		SendMessage(bot, callback.Message.Chat.ID, "Unknown request")
		return
	}

	// Create a fake message object to reuse ReplyForAMessage
	fakeMessage := &tgbotapi.Message{
		Chat: callback.Message.Chat,
		Text: question,
	}

	ReplyForAMessage(bot, fakeMessage)
}

func Listen(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Handle callback queries (button presses)
		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
			continue
		}

		// Handle regular messages
		if update.Message == nil {
			continue
		}

		// Handle /start command
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			sendStartMessage(bot, update.Message.Chat.ID)
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

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	Listen(bot)
}
