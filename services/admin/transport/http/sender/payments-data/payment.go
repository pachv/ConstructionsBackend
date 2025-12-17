package paymentsdata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PaymentData соответствует структуре платежа
type PaymentData struct {
	Id            string `json:"id"`
	UserId        string `json:"userId"`
	PriceType     string `json:"priceType"`
	PriceAmount   string `json:"priceAmount"`
	RevardType    string `json:"revardType"`
	RevardAmount  string `json:"revardAmount"`
	PaymentStatus string `json:"paymentStatus"`
	PaymentTime   string `json:"paymentTime"`
}

// PaymentsResponse - структура ответа от микросервиса
type PaymentsResponse struct {
	Payments   []*PaymentData `json:"payments"`
	PageAmount int            `json:"pageAmount"`
}

// FetchPaymentsData делает GET-запрос на микросервис оплат
func FetchPaymentsData(page int, search, orderBy string) (*PaymentsResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("http://payment:8080/admin/payments?page=%d&search=%s&orderBy=%s",
		page, search, orderBy)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call payment service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}

	var paymentsResp PaymentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentsResp); err != nil {
		return nil, fmt.Errorf("failed to decode payment service response: %w", err)
	}

	return &paymentsResp, nil
}
