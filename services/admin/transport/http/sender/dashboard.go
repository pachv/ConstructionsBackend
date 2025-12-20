package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const constructionsBaseURL = "http://constructions_service:8080"
const constructionsDashboardPath = "/admin/dashboard" // <-- если у тебя другой роут, поменяй тут

type ConstructionsDashboardStats struct {
	TotalUsers      int `json:"totalUsers"`
	TotalOrders     int `json:"totalOrders"`
	OrdersToday     int `json:"ordersToday"`
	OrdersLastMonth int `json:"ordersLastMonth"`
}

func GetConstructionsDashboardStats(ctx context.Context) (*ConstructionsDashboardStats, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, constructionsBaseURL+constructionsDashboardPath, nil)
	if err != nil {
		return nil, fmt.Errorf("dashboard: create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dashboard: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dashboard: bad status %d", resp.StatusCode)
	}

	var out ConstructionsDashboardStats
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("dashboard: decode: %w", err)
	}

	return &out, nil
}
