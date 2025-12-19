package sender

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ListResponse struct {
	Items []SectionShort `json:"items"`
}
type SectionShort struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	HasGallery bool   `json:"hasGallery"`
	HasCatalog bool   `json:"hasCatalog"`
}

type SectionFull struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	HasGallery bool   `json:"hasGallery"`
	HasCatalog bool   `json:"hasCatalog"`

	Gallery []GalleryItem `json:"gallery"`
	Catalog *Catalog      `json:"catalog"`
}

type GalleryItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	SortOrder int    `json:"sortOrder"`
}

type Catalog struct {
	Categories []CatalogCategory `json:"categories"`
	Items      []CatalogItem     `json:"items"`
}

type CatalogCategory struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	SortOrder int    `json:"sortOrder"`
}

type CatalogItem struct {
	ID         string   `json:"id"`
	CategoryId string   `json:"categoryId"`
	Title      string   `json:"title"`
	PriceRub   int      `json:"priceRub"`
	ImageUrl   string   `json:"imageUrl"`
	Badges     []string `json:"badges"`
	Specs      []Spec   `json:"specs"`
	SortOrder  int      `json:"sortOrder"`
}

type Spec struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Client struct {
	BaseURL string
	Http    *http.Client
}

func NewSectionsSender(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetList() (ListResponse, error) {
	var out ListResponse
	req, _ := http.NewRequest(http.MethodGet, c.BaseURL+"/api/v1/sections", nil)
	resp, err := c.Http.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return out, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (c *Client) GetBySlug(slug string) (SectionFull, error) {
	var out SectionFull
	req, _ := http.NewRequest(http.MethodGet, c.BaseURL+"/api/v1/sections/"+slug, nil)
	resp, err := c.Http.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return out, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return out, json.NewDecoder(resp.Body).Decode(&out)
}
