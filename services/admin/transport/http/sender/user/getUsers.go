package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UsersResponse struct {
	Users      []UserData `json:"users"`
	PageAmount int        `json:"pageAmount"`
}

type UserData struct {
	Id              string    `db:"id" json:"id"`
	Username        string    `db:"username" json:"username"`
	TelegramId      int64     `db:"telegram_id" json:"telegramId"`
	FirstLogin      time.Time `db:"time_first_login" json:"firstLogin"`
	LastAuth        time.Time `db:"last_auth" json:"lastAuth"`
	PremiumDaysLeft int64     `db:"premium_days_left" json:"premiumDaysLeft"`
	RequestsLeft    int64     `db:"requestsLeft" json:"requestsLeft"`
	Promo           string    `db:"promo" json:"promo"`
	ReferalSenderId string    `db:"referal_telegram_id" json:"referalSenderId"`
}

// FetchUsersData делает GET запрос к микросервису user-service по адресу
// http://user:8080/admin/users-data?page=&search=&orderBy=
func FetchUsersData(page int, search, orderBy string) (*UsersResponse, error) {
	baseURL := "http://user:8080/admin/users-data"

	// Формируем query параметры
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	if search != "" {
		params.Add("search", search)
	}
	if orderBy != "" {
		params.Add("orderBy", orderBy)
	}

	// Формируем итоговый URL
	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Делаем запрос
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %s", resp.Status)
	}

	// Декодируем JSON
	var data UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}
