package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DashboardUserData struct {
	ActiveSubscriptionAmount int `json:"ActiveSubscriptionAmount"`
	TotalUsers               int `json:"totalUsers"`
}

func GetDashboardUserData() (*DashboardUserData, error) {
	url := "http://user:8080/admin/dashboard-data"

	// Делаем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при GET запросе: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %w", err)
	}

	// Парсим JSON
	var data DashboardUserData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге JSON: %w", err)
	}

	return &data, nil
}
