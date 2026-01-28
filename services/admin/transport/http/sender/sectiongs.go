package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// БАЗОВЫЕ URL — сделай как у тебя принято (через config/const).
// Я даю явные константы, чтобы было 1:1 как GetAdminEmail().
const (
	PublicAPIBaseURL             = "http://constructions_service:8080"
	GetAdminSectionsAllURL       = PublicAPIBaseURL + "/admin/sections/all"
	GetAdminSectionFullBySlugURL = PublicAPIBaseURL + "/admin/sections/%s/full"
	CreateSectionFormURL         = PublicAPIBaseURL + "/admin/sections/create-form"
	DeleteSectionURL             = PublicAPIBaseURL + "/admin/sections/%s"
)

type SectionsListResponse struct {
	Items []SectionSummary `json:"items"`
}

type SectionSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Label      string `json:"label"`
	Slug       string `json:"slug"`
	Image      string `json:"image"`
	HasGallery bool   `json:"hasGallery"`
	HasCatalog bool   `json:"hasCatalog"`
}

type SectionFull struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Label           string        `json:"label"`
	Slug            string        `json:"slug"`
	Image           string        `json:"image"`
	AdvantegesText  string        `json:"advantegesText"`
	AdvantegesArray []string      `json:"advantegesArray"`
	HasGallery      bool          `json:"hasGallery"`
	HasCatalog      bool          `json:"hasCatalog"`
	Gallery         []GalleryItem `json:"gallery"`
	Catalog         *Catalog      `json:"catalog"`
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
	ID         string `json:"id"`
	CategoryId string `json:"categoryId"`
	Title      string `json:"title"`
	PriceRub   int    `json:"priceRub"`
	ImageUrl   string `json:"imageUrl"`
	SortOrder  int    `json:"sortOrder"`
	Specs      []Spec `json:"specs"`
}

type Spec struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// =========================
// GET list
// =========================

func GetAdminSectionsAll(ctx context.Context) (SectionsListResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, GetAdminSectionsAllURL, nil)
	if err != nil {
		return SectionsListResponse{}, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return SectionsListResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return SectionsListResponse{}, fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	var out SectionsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return SectionsListResponse{}, err
	}
	if out.Items == nil {
		out.Items = []SectionSummary{}
	}
	return out, nil
}

// =========================
// GET full by slug
// =========================

func GetAdminSectionFullBySlug(ctx context.Context, slug string) (SectionFull, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return SectionFull{}, fmt.Errorf("slug is required")
	}

	u := fmt.Sprintf(GetAdminSectionFullBySlugURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return SectionFull{}, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return SectionFull{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return SectionFull{}, fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	var out SectionFull
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return SectionFull{}, err
	}

	// чтобы шаблон не падал
	if out.AdvantegesArray == nil {
		out.AdvantegesArray = []string{}
	}
	if out.Gallery == nil {
		out.Gallery = []GalleryItem{}
	}
	if out.Catalog == nil {
		out.Catalog = &Catalog{Categories: []CatalogCategory{}, Items: []CatalogItem{}}
	}
	if out.Catalog.Categories == nil {
		out.Catalog.Categories = []CatalogCategory{}
	}
	if out.Catalog.Items == nil {
		out.Catalog.Items = []CatalogItem{}
	}

	return out, nil
}

// =========================
// CREATE basic (multipart)
// =========================

func CreateSectionBasic(ctx context.Context, title, advText string, advs []string, imagePath string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(imagePath) == "" {
		return fmt.Errorf("imagePath is required")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("title", title)
	_ = w.WriteField("advantegesText", strings.TrimSpace(advText))

	for _, a := range advs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		_ = w.WriteField("advanteges[]", a)
	}

	f, err := os.Open(imagePath)
	if err != nil {
		_ = w.Close()
		return err
	}
	defer f.Close()

	part, err := w.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		_ = w.Close()
		return err
	}
	if _, err := io.Copy(part, f); err != nil {
		_ = w.Close()
		return err
	}
	_ = w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, CreateSectionFormURL, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	return nil
}

// =========================
// DELETE
// =========================

func DeleteSection(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	u := fmt.Sprintf(DeleteSectionURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}
	return nil
}
