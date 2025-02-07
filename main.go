package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	token      = "YOUR_OAUTH_TOKEN"
	counterID  = "YOUR_COUNTER_ID"
	metricsURL = "https://api-metrika.yandex.net/stat/v1/data"
)

type YandexMetrikaResponse struct {
	Data []struct {
		Metrics []float64 `json:"metrics"`
	} `json:"data"`
}

func getMetrikaStats() {
	params := url.Values{}
	params.Set("ids", counterID)
	params.Set("metrics", "ym:s:visits,ym:s:users,ym:s:pageviews")
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

	fmt.Println("Статистика:")
	for _, data := range result.Data {
		fmt.Printf("Визиты: %.0f, Пользователи: %.0f, Просмотры: %.0f\n",
			data.Metrics[0], data.Metrics[1], data.Metrics[2])
	}
}

func main() {
	getMetrikaStats()
}
