package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// const constructionsBaseURL = "http://constructions_service:8080"

/*
=========================
 DTO (local entity)
=========================
*/

type SiteSectionSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Label      string `json:"label"`
	Slug       string `json:"slug"`
	Image      string `json:"image"`
	HasGallery bool   `json:"hasGallery"`
	HasCatalog bool   `json:"hasCatalog"`
}

type SiteSectionGallery struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	SortOrder int    `json:"sortOrder"`
}

type SiteSection struct {
	ID              string               `json:"id"`
	Title           string               `json:"title"`
	Label           string               `json:"label"`
	Slug            string               `json:"slug"`
	Image           string               `json:"image"`
	AdvantegesText  string               `json:"advantegesText"`
	AdvantegesArray []string             `json:"advantegesArray"`
	HasGallery      bool                 `json:"hasGallery"`
	HasCatalog      bool                 `json:"hasCatalog"`
	Gallery         []SiteSectionGallery `json:"gallery,omitempty"`
}

type AdminSectionsPage struct {
	Items      []*SiteSectionSummary `json:"items"`
	Page       int                   `json:"page"`
	PageAmount int                   `json:"pageAmount"`
	Total      int                   `json:"total"`
}

/*
=========================
 HTTP helpers
=========================
*/

func doRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	return b, resp.StatusCode, nil
}

/*
=========================
 API
=========================
*/

// GET /admin/sections
func GetAdminSections(
	ctx context.Context,
	page int,
	search string,
	orderBy string,
) (*AdminSectionsPage, error) {

	u, err := url.Parse(constructionsBaseURL + "/admin/sections")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if search != "" {
		q.Set("search", search)
	}
	if orderBy != "" {
		q.Set("orderBy", orderBy)
	}
	u.RawQuery = q.Encode()

	body, status, err := doRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("get admin sections failed: status=%d body=%s", status, string(body))
	}

	var out AdminSectionsPage
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// GET /admin/sections/:slug
func GetAdminSectionBySlug(
	ctx context.Context,
	slug string,
) (*SiteSection, error) {

	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}

	url := constructionsBaseURL + "/admin/sections/" + url.PathEscape(slug)

	body, status, err := doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("get admin section failed: status=%d body=%s", status, string(body))
	}

	var out SiteSection
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
