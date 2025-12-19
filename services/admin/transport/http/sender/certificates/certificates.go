package certificates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const mainServiceURL = "http://constructions_service:8080/admin"

type Item struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
}

func GetAll(ctx context.Context) ([]Item, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mainServiceURL+"/certificates", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("certificates.GetAll: unexpected status %d", resp.StatusCode)
	}

	var out []Item
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}
