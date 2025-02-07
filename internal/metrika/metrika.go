package metrika

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
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

func GetStats() {
	LoadEnv()

	token := os.Getenv("YANDEX_METRIKA_TOKEN")
	counterID := os.Getenv("YANDEX_METRIKA_COUNTER_ID")

	if token == "" || counterID == "" {
		log.Fatal("Error: YANDEX_METRIKA_TOKEN or YANDEX_METRIKA_COUNTER_ID is not set")
	}

	params := url.Values{}
	params.Set("ids", counterID)
	params.Set("metrics", "ym:s:pageviews,ym:s:visits,ym:s:users")
	params.Set("date1", "today")
	params.Set("date2", "today")
	params.Set("accuracy", "full")

	req, err := http.NewRequest("GET", metricsURL+"?"+params.Encode(), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: Response status %d", resp.StatusCode)
	}

	var result YandexMetrikaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	fmt.Printf("Статистика за %s:\n", today)

	if len(result.Data) > 0 {
		data := result.Data[0]
		fmt.Printf("Просмотры: %.0f, Визиты: %.0f, Посетители: %.0f\n",
			data.Metrics[0], data.Metrics[1], data.Metrics[2])
	} else {
		fmt.Println("No data available for today")
	}
}
