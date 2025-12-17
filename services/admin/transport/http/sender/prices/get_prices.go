package prices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Price struct {
	PriceId           string  `json:"PriceId"`
	LanguageId        string  `json:"LanguageId"`
	PriceTitle        string  `json:"PriceTitle"`
	PriceDescription  string  `json:"PriceDescription"`
	PriceRewardType   string  `json:"PriceRewardType"`
	PriceRewardAmount int     `json:"PriceRewardAmount"`
	PriceTonId        string  `json:"PriceTonId"`
	PriceStartsId     string  `json:"PriceStartsId"`
	PriceTonAmount    float64 `json:"PriceTonAmount"`
	PriceStarsAmount  int     `json:"PriceStarsAmount"`
}

// структура всего ответа
type pricesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Prices []Price `json:"prices"`
	} `json:"data"`
}

// функция получения цен
func GetPrices() ([]Price, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("http://payment:8080/admin/prices")
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var result pricesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data.Prices, nil
}

func main() {
	prices, err := GetPrices()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("✅ Received %d prices\n\n", len(prices))
	for _, p := range prices {
		fmt.Printf("%s [%s]: %.4f TON, %d Stars (lang: %s)\n",
			p.PriceTitle, p.PriceRewardType, p.PriceTonAmount, p.PriceStarsAmount, p.LanguageId)
	}
}
