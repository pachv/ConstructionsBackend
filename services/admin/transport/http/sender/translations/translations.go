package translations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TranslationsResponse struct {
	Data struct {
		Translations string `json:"translations"`
	} `json:"data"`
	Status string `json:"status"`
}

// Функция для получения переводов
func GetTranslations() (translations string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://settings:8080/admin/get-translations", nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result TranslationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "nil", fmt.Errorf("decode response: %w", err)
	}

	return result.Data.Translations, nil
}

type SetTransaltionsRequest struct {
	Translation string `json:"translations"`
}

// структура для возможного ответа
type SetTranslationsResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// функция для отправки перевода
func SetTranslations(translation string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqBody := SetTransaltionsRequest{
		Translation: translation,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://settings:8080/admin/set-translations", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result SetTranslationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
