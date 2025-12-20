package certificates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const mainServiceURL = "http://constructions_service:8080/admin"

type Item struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
}

type ListResponse struct {
	Items      []Item `json:"items"`
	Page       int    `json:"page"`
	PageAmount int    `json:"pageAmount"`
}

func GetAll(ctx context.Context, page int, search string) (*ListResponse, error) {
	if page < 1 {
		page = 1
	}

	u, err := url.Parse(mainServiceURL + "/certificates")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("certificates.GetAll: unexpected status %d", resp.StatusCode)
	}

	var out ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
