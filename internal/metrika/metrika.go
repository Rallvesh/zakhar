package metrika

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/rallvesh/zakhar/internal/logger"
)

const metricsURL = "https://api-metrika.yandex.net/stat/v1/data"

type YandexMetrikaResponse struct {
	Data []struct {
		Dimensions []struct {
			Name string `json:"name"`
		} `json:"dimensions"`
		Metrics []float64 `json:"metrics"`
	} `json:"data"`
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

func GetUserStats() string {
	LoadEnv()

	generalStats, err := FetchStats("ym:s:pageviews,ym:s:visits,ym:s:users", "")
	if err != nil || len(generalStats.Data) == 0 {
		return "Ошибка получения статистики по пользователям."
	}

	return fmt.Sprintf(
		"👥 Пользователи:\n"+
			"- Просмотры: %.0f\n"+
			"- Визиты: %.0f\n"+
			"- Пользователи: %.0f",
		generalStats.Data[0].Metrics[0], // Просмотры
		generalStats.Data[0].Metrics[1], // Визиты
		generalStats.Data[0].Metrics[2], // Пользователи
	)
}

func GetTrafficStats() string {
	LoadEnv()

	trafficStats, err := FetchStats("ym:s:visits", "ym:s:trafficSource")
	if err != nil || len(trafficStats.Data) == 0 {
		return "Ошибка получения статистики по источникам трафика."
	}

	trafficSources := ""
	for _, row := range trafficStats.Data {
		if len(row.Dimensions) > 0 {
			trafficSources += fmt.Sprintf("- %s: %.0f\n", row.Dimensions[0].Name, row.Metrics[0])
		}
	}
	if trafficSources == "" {
		trafficSources = "Нет данных по источникам трафика."
	}

	return fmt.Sprintf(
		"🚦 Источники трафика:\n%s",
		trafficSources,
	)
}

func FetchStats(metrics string, dimensions string) (YandexMetrikaResponse, error) {
	token := os.Getenv("YANDEX_METRIKA_TOKEN")
	counterID := os.Getenv("YANDEX_METRIKA_COUNTER_ID")

	logger := logger.Init()

	if token == "" || counterID == "" {
		logger.Error("YANDEX_METRIKA_TOKEN or YANDEX_METRIKA_COUNTER_ID is not set")
	}

	params := url.Values{}
	params.Set("ids", counterID)
	params.Set("metrics", metrics)
	params.Set("date1", "today")
	params.Set("date2", "today")
	params.Set("accuracy", "full")
	if dimensions != "" {
		params.Set("dimensions", dimensions)
	}

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
		logger.Error("Response status", slog.Int("status_code", resp.StatusCode))
	}

	var result YandexMetrikaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("Error parsing response", slog.Any("error", err))
	}

	return result, nil
}
