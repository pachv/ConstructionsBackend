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

const (
	PublicAPIBaseURL             = "http://constructions_service:8080"
	GetAdminSectionsAllURL       = PublicAPIBaseURL + "/admin/sections/all"
	GetAdminSectionFullBySlugURL = PublicAPIBaseURL + "/admin/sections/%s/full"
	CreateSectionFormURL         = PublicAPIBaseURL + "/admin/sections/create-form"
	DeleteSectionURL             = PublicAPIBaseURL + "/admin/sections/%s"
	UpdateSectionURL             = PublicAPIBaseURL + "/admin/sections/%s"
	ToggleGalleryURL             = PublicAPIBaseURL + "/admin/sections/%s/gallery/toggle"
	ToggleCatalogURL             = PublicAPIBaseURL + "/admin/sections/%s/catalog/toggle"
	AddGalleryPhotoURL           = PublicAPIBaseURL + "/admin/sections/%s/gallery"
	DeleteGalleryPhotoURL        = PublicAPIBaseURL + "/admin/sections/%s/gallery/%s"
	UploadGalleryImageURL        = PublicAPIBaseURL + "/admin/sections/%s/gallery/upload"
	AddCatalogCategoryURL        = PublicAPIBaseURL + "/admin/sections/%s/catalog/categories"
	DeleteCatalogCategoryURL     = PublicAPIBaseURL + "/admin/sections/%s/catalog/categories/%s"
	UpdateCatalogURL             = PublicAPIBaseURL + "/admin/sections/%s/catalog"
	AddCatalogItemURL            = PublicAPIBaseURL + "/admin/sections/%s/catalog/items"
	DeleteCatalogItemURL         = PublicAPIBaseURL + "/admin/sections/%s/catalog/items/%s"
	UploadCatalogItemImageURL    = PublicAPIBaseURL + "/admin/sections/%s/catalog/items/upload"
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
	ID         string   `json:"id"`
	CategoryId string   `json:"categoryId"`
	Title      string   `json:"title"`
	PriceRub   int      `json:"priceRub"`
	ImageUrl   string   `json:"imageUrl"`
	SortOrder  int      `json:"sortOrder"`
	Badges     []string `json:"badges"`
	Specs      []Spec   `json:"specs"`
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
// UPDATE section (multipart)
// =========================

type UpdateSectionInput struct {
	Title           string
	Slug            string
	AdvantegesText  string
	AdvantegesArray []string
	HasGallery      bool
	HasCatalog      bool
	ImagePath       string // опционально
}

func UpdateSection(ctx context.Context, id string, input UpdateSectionInput) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("title", input.Title)
	_ = w.WriteField("slug", input.Slug)
	_ = w.WriteField("advantegesText", input.AdvantegesText)
	_ = w.WriteField("hasGallery", fmt.Sprintf("%t", input.HasGallery))
	_ = w.WriteField("hasCatalog", fmt.Sprintf("%t", input.HasCatalog))

	advJSON, _ := json.Marshal(input.AdvantegesArray)
	_ = w.WriteField("advantegesArray", string(advJSON))

	if input.ImagePath != "" {
		f, err := os.Open(input.ImagePath)
		if err == nil {
			defer f.Close()
			part, err := w.CreateFormFile("image", filepath.Base(input.ImagePath))
			if err == nil {
				_, _ = io.Copy(part, f)
			}
		}
	}

	_ = w.Close()

	u := fmt.Sprintf(UpdateSectionURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, &buf)
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

	if resp.StatusCode != http.StatusOK {
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

// =========================
// TOGGLE Gallery/Catalog
// =========================

func ToggleSectionGallery(ctx context.Context, id string, enabled bool) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	payload := map[string]bool{"enabled": enabled}
	body, _ := json.Marshal(payload)

	u := fmt.Sprintf(ToggleGalleryURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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

func ToggleSectionCatalog(ctx context.Context, id string, enabled bool) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	payload := map[string]bool{"enabled": enabled}
	body, _ := json.Marshal(payload)

	u := fmt.Sprintf(ToggleCatalogURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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

// =========================
// GALLERY
// =========================

type AddGalleryPhotoInput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	SortOrder int    `json:"sortOrder"`
}

func AddGalleryPhoto(ctx context.Context, slug string, input AddGalleryPhotoInput) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("slug is required")
	}

	body, _ := json.Marshal(input)
	u := fmt.Sprintf(AddGalleryPhotoURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}
	return nil
}

func DeleteGalleryPhoto(ctx context.Context, slug, photoID string) error {
	slug = strings.TrimSpace(slug)
	photoID = strings.TrimSpace(photoID)
	if slug == "" || photoID == "" {
		return fmt.Errorf("slug and photoID are required")
	}

	u := fmt.Sprintf(DeleteGalleryPhotoURL, slug, photoID)
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

func UploadGalleryImage(ctx context.Context, slug, imagePath string) (string, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", fmt.Errorf("slug is required")
	}
	if imagePath == "" {
		return "", fmt.Errorf("imagePath is required")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	part, err := w.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, f); err != nil {
		return "", err
	}
	_ = w.Close()

	u := fmt.Sprintf(UploadGalleryImageURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	if err := json.Unmarshal(body, &result); err == nil {
		if url, ok := result["url"]; ok {
			return url, nil
		}
	}
	return string(body), nil
}

// =========================
// CATALOG - Categories
// =========================

type AddCatalogCategoryInput struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	SortOrder int    `json:"sortOrder"`
}

func AddCatalogCategory(ctx context.Context, slug string, input AddCatalogCategoryInput) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("slug is required")
	}

	body, _ := json.Marshal(input)
	u := fmt.Sprintf(AddCatalogCategoryURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}
	return nil
}

func DeleteCatalogCategory(ctx context.Context, slug, categoryID string) error {
	slug = strings.TrimSpace(slug)
	categoryID = strings.TrimSpace(categoryID)
	if slug == "" || categoryID == "" {
		return fmt.Errorf("slug and categoryID are required")
	}

	u := fmt.Sprintf(DeleteCatalogCategoryURL, slug, categoryID)
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

func UpdateCatalog(ctx context.Context, slug string, catalog *Catalog) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("slug is required")
	}

	body, _ := json.Marshal(catalog)
	u := fmt.Sprintf(UpdateCatalogURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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

// =========================
// CATALOG - Items
// =========================

type AddCatalogItemInput struct {
	ID         string   `json:"id"`
	CategoryId string   `json:"categoryId"`
	Title      string   `json:"title"`
	PriceRub   int      `json:"priceRub"`
	ImageUrl   string   `json:"imageUrl"`
	SortOrder  int      `json:"sortOrder"`
	Badges     []string `json:"badges"`
	Specs      []Spec   `json:"specs"`
}

func AddCatalogItem(ctx context.Context, slug string, input AddCatalogItemInput) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("slug is required")
	}

	body, _ := json.Marshal(input)
	u := fmt.Sprintf(AddCatalogItemURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}
	return nil
}

func DeleteCatalogItem(ctx context.Context, slug, itemID string) error {
	slug = strings.TrimSpace(slug)
	itemID = strings.TrimSpace(itemID)
	if slug == "" || itemID == "" {
		return fmt.Errorf("slug and itemID are required")
	}

	u := fmt.Sprintf(DeleteCatalogItemURL, slug, itemID)
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

func UploadCatalogItemImage(ctx context.Context, slug, imagePath string) (string, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", fmt.Errorf("slug is required")
	}
	if imagePath == "" {
		return "", fmt.Errorf("imagePath is required")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	part, err := w.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, f); err != nil {
		return "", err
	}
	_ = w.Close()

	u := fmt.Sprintf(UploadCatalogItemImageURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	if err := json.Unmarshal(body, &result); err == nil {
		if url, ok := result["url"]; ok {
			return url, nil
		}
	}
	return string(body), nil
}
