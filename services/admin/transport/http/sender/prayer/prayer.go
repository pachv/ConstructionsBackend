package prayer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Prayer struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	EngName     string `json:"EngName"`
	Description string `json:"Description"`
	LanguageId  string `json:"LanguageId"`
}

func GetPrayers() ([]Prayer, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://prayer:8080/admin/get-prayers")
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// структура для парсинга полного ответа
	var respData struct {
		Prayers []Prayer `json:"prayers"`
	}

	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	// возвращаем только массив prayers
	return respData.Prayers, nil
}

type PrayerUpdate struct {
	PrayerID          string `json:"prayerId"`
	PrayerName        string `json:"prayerName"`
	PrayerDescription string `json:"prayerDescription"`
}

// структура для всего запроса
type SetPrayersRequest struct {
	Prayers []PrayerUpdate `json:"prayers"`
}

func SetPrayers(prayers []PrayerUpdate) error {
	// сериализация данных в JSON
	reqBody := SetPrayersRequest{Prayers: prayers}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// делаем POST-запрос
	resp, err := client.Post("http://prayer:8080/admin/set-prayers", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
