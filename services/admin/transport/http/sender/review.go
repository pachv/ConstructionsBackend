package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// const constructionsBaseURL = "http://constructions_service:8080"
const constructionsAdminReviewsPath = "/admin/reviews"

type AdminReview struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Position   string    `json:"position"`
	Text       string    `json:"text"`
	Rating     int       `json:"rating"`
	ImagePath  string    `json:"imagePath"`
	Consent    bool      `json:"consent"`
	CanPublish bool      `json:"canPublish"`
	CreatedAt  time.Time `json:"createdAt"`
}

type GetReviewsResponse struct {
	Items      []AdminReview `json:"items"`
	Page       int           `json:"page"`
	PageAmount int           `json:"pageAmount"`
}

func GetAdminReviews(ctx context.Context, page int, search, orderBy string) (*GetReviewsResponse, error) {
	if page < 1 {
		page = 1
	}
	search = strings.TrimSpace(search)
	orderBy = strings.TrimSpace(orderBy)

	u, err := url.Parse(constructionsBaseURL + constructionsAdminReviewsPath)
	if err != nil {
		return nil, fmt.Errorf("reviews: parse url: %w", err)
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}
	if orderBy != "" {
		q.Set("orderBy", orderBy)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("reviews: create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("reviews: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("reviews: bad status %d", resp.StatusCode)
	}

	var out GetReviewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("reviews: decode: %w", err)
	}

	return &out, nil
}
