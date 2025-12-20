package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type getAdminEmailResponse struct {
	Status string `json:"status"`
	Email  string `json:"email"`
}

func GetAdminEmail(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		GetAdminEmailURL,
		nil,
	)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	var out getAdminEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Email, nil
}
