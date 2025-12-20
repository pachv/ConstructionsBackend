package gallery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const publicAPI = "http://constructions_service:8080/api/v1"

type Category struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Photo struct {
	ID           string `json:"id"`
	CategorySlug string `json:"category_slug"`
	Alt          string `json:"alt"`
	Image        string `json:"image"` // это image_path = только имя файла
	SortOrder    int    `json:"sort_order"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type categoriesResp struct {
	Items []Category `json:"items"`
}

type photosResp struct {
	Items []Photo `json:"items"`
}

func GetCategories(ctx context.Context) ([]Category, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, publicAPI+"/gallery/categories", nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 10 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("GetCategories bad status: %s", res.Status)
	}

	var out categoriesResp
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func GetPhotosBySlug(ctx context.Context, slug string) ([]Photo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, publicAPI+"/gallery/"+slug+"/photos", nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 10 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("GetPhotosBySlug bad status: %s", res.Status)
	}

	var out photosResp
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Items, nil
}
