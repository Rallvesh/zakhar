package bot

import (
	"log"
	"log/slog"
	"os"

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

			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, "+update.Message.From.FirstName+"!")
				bot.Send(msg)
				// logger.Info("/start", slog.String("user", update.Message.From.UserName), slog.Int64("chat_id", update.Message.Chat.ID))
			case "stats":
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
