package metrika

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rallvesh/zakhar/internal/logger"
)

const metricsURL = "https://api-metrika.yandex.net/stat/v1/data"

type YandexMetrikaResponse struct {
	Data []struct {
		Metrics []float64 `json:"metrics"`
	} `json:"data"`
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

func GetStats() string {
	LoadEnv()

	logger := logger.Init()

	token := os.Getenv("YANDEX_METRIKA_TOKEN")
	counterID := os.Getenv("YANDEX_METRIKA_COUNTER_ID")

	if token == "" || counterID == "" {
		logger.Error("YANDEX_METRIKA_TOKEN or YANDEX_METRIKA_COUNTER_ID is not set")
	}

	params := url.Values{}
	params.Set("ids", counterID)
	params.Set("metrics", "ym:s:pageviews,ym:s:visits,ym:s:users")
	params.Set("date1", "today")
	params.Set("date2", "today")
	params.Set("accuracy", "full")

	req, err := http.NewRequest("GET", metricsURL+"?"+params.Encode(), nil)
	if err != nil {
		logger.Error("Error creating request", slog.Any("error", err))
	}
	req.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending request", slog.Any("error", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Response status", slog.Any("error", resp.StatusCode))
	}

	var result YandexMetrikaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("Error parsing response", slog.Any("error", err))
	}

	today := time.Now().Format("2006-01-02")

	if len(result.Data) > 0 {
		data := result.Data[0]
		return fmt.Sprintf("📊 Статистика за %s:\nВизиты: %.0f\nПользователи: %.0f\nПросмотры: %.0f",
			today, data.Metrics[0], data.Metrics[1], data.Metrics[2])
	}
	return "Нет данных за сегодняшний день"
}
