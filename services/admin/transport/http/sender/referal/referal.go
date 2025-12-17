package referal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const referalBaseURL = "http://referal:8080/inside/admin/bonuses"

// --- СТРУКТУРЫ ---

type Bonuses struct {
	TotalAmount     float64 `json:"totalAmount"`
	PurchasedAmount float64 `json:"purchasedAmount"`
	TotalGift       float64 `json:"totalGift"`
	PurchasedGift   float64 `json:"purchasedGift"`
	URL             string  `json:"url"`
	EnglishText     string  `json:"englishText"`
	ArabText        string  `json:"arabText"`
}

type okResponse struct {
	Status string  `json:"status"`
	Data   Bonuses `json:"data"`
}

// --- ФУНКЦИИ ПОЛУЧЕНИЯ ---

// GetReferalData получает все данные (включая url, englishText, arabText)
func GetReferalData() (*Bonuses, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, referalBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var result okResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("server returned status: %s", result.Status)
	}

	return &result.Data, nil
}

// --- ФУНКЦИИ ОТПРАВКИ ---

// SetReferalData отправляет обновлённые данные обратно в сервис referal
func SetReferalData(data *Bonuses) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка маршала JSON: %w", err)
	}

	resp, err := http.Post(referalBaseURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул статус: %v", resp.Status)
	}

	return nil
}
