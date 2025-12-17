package ton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Структуры для парсинга JSON ответа
type TonDataResponse struct {
	Data   TonData `json:"data"`
	Status string  `json:"status"`
}

type TonData struct {
	URL    string `json:"url"`
	Wallet string `json:"wallet"`
}

// GetTonData делает GET-запрос и возвращает данные Ton
func GetTonData() (TonData, error) {
	url := "http://payment:8080/admin/ton-data"

	resp, err := http.Get(url)
	if err != nil {
		return TonData{}, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TonData{}, fmt.Errorf("сервер вернул статус: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TonData{}, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var result TonDataResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return TonData{}, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return result.Data, nil
}

type TonDataRequest struct {
	Wallet string `json:"wallet"`
	URL    string `json:"url"`
}

func SendTonData(wallet, urlStr string) (TonDataResponse, error) {
	requestBody := TonDataRequest{
		Wallet: wallet,
		URL:    urlStr,
	}

	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		return TonDataResponse{}, fmt.Errorf("ошибка маршала JSON: %w", err)
	}

	resp, err := http.Post("http://payment:8080/admin/ton-data", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return TonDataResponse{}, fmt.Errorf("ошибка отправки POST-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TonDataResponse{}, fmt.Errorf("сервер вернул статус: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TonDataResponse{}, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var result TonDataResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return TonDataResponse{}, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return result, nil
}
