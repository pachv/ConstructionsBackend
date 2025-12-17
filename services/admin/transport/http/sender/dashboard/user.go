package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type dashboardPaymentsResponse struct {
	Data struct {
		PremiumSold  int `json:"premiumSold"`
		RequestsSold int `json:"requestsSold"`
	} `json:"data"`
	Status string `json:"status"`
}

// структура, которую будет возвращать функция
type DashboardPaymentsData struct {
	PremiumSold  int
	RequestsSold int
}

// функция для получения данных
func GetDashboardPayments() (*DashboardPaymentsData, error) {
	url := "http://payment:8080/admin/dashboard"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при GET запросе: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %s", resp.Status)
	}

	var result dashboardPaymentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге JSON: %w", err)
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("сервер вернул статус: %s", result.Status)
	}

	data := DashboardPaymentsData{
		PremiumSold:  result.Data.PremiumSold,
		RequestsSold: result.Data.RequestsSold,
	}

	return &data, nil
}
