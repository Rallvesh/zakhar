package bot

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rallvesh/zakhar/internal/logger"
	"github.com/rallvesh/zakhar/internal/metrika"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

func Start() {
	LoadEnv()

	logger := logger.Init()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		logger.Error("TELEGRAM_BOT_TOKEN is not set")
	}

	// Get allowed chat ID and user IDs from environment variables
	allowedChatID, _ := strconv.ParseInt(os.Getenv("ALLOWED_CHAT_ID"), 10, 64)
	allowedUsers := parseAllowedUsers(os.Getenv("ALLOWED_USER_IDS"))

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Error creating bot", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Bot is started", slog.String("name", bot.Self.FirstName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			logger.Info(update.Message.Text, slog.String("user", update.Message.From.UserName), slog.Int64("chat_id", update.Message.Chat.ID))

			userID := update.Message.From.ID
			chatID := update.Message.Chat.ID

			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, "+update.Message.From.FirstName+"!")
				bot.Send(msg)
				// logger.Info("/start", slog.String("user", update.Message.From.UserName), slog.Int64("chat_id", update.Message.Chat.ID))
			case "stats":
				// Check if the chat or user is allowed
				if chatID != allowedChatID && !allowedUsers[userID] {
					logger.Warn("Unauthorized access attempt", slog.String("user", update.Message.From.UserName), slog.Int64("chat_id", chatID))
					msg := tgbotapi.NewMessage(chatID, "You are not authorized to use this command")
					bot.Send(msg)
					continue
				}

				stats := metrika.GetStats()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, stats)
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
				bot.Send(msg)
			}
		}
	}
}

// parseAllowedUsers parses a comma-separated list of user IDs from environment variables
func parseAllowedUsers(envVar string) map[int64]bool {
	users := make(map[int64]bool)
	userIDs := strings.Split(envVar, ",")
	for _, idStr := range userIDs {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err == nil {
			users[id] = true
		}
	}
	return users
}
