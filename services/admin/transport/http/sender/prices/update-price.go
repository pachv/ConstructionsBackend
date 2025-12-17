package prices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdatePrice struct {
	PriceId           string  `json:"PriceId"`
	PriceTitle        string  `json:"PriceTitle"`
	PriceDescription  string  `json:"PriceDescription"`
	PriceRewardAmount int     `json:"PriceRewardAmount"`
	PriceTonId        string  `json:"PriceTonId"`
	PriceStartsId     string  `json:"PriceStartsId"`
	PriceTonAmount    float64 `json:"PriceTonAmount"`
	PriceStarsAmount  int     `json:"PriceStarsAmount"`
}

type PricesRequest struct {
	Prices []UpdatePrice `json:"prices"`
}

// UpdatePrices отправляет массив цен на http://payment:8080/admin/prices
func UpdatePrices(prices []UpdatePrice) error {
	url := "http://payment:8080/admin/prices"

	body := PricesRequest{Prices: prices}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул статус: %s", resp.Status)
	}

	fmt.Println("Цены успешно обновлены")
	return nil
}
