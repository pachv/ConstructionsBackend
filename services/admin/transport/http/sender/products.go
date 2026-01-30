package sender

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// ✅ для списков в pages (GET)
	ConstructionsBaseURL = "http://constructions_service:8080"

	// ✅ для CRUD через admin-service (handlers)
	AdminServiceBaseURL = "http://localhost:80/admin-service"

	client = &http.Client{Timeout: 7 * time.Second}
)

func joinConstructions(path string) string {
	return strings.TrimRight(ConstructionsBaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}
func joinAdminService(path string) string {
	return strings.TrimRight(AdminServiceBaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func doJSON[T any](ctx context.Context, method, fullURL string, in any) (T, error) {
	var zero T

	var body *bytes.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return zero, err
		}
		body = bytes.NewReader(b)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return zero, err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return zero, sql.ErrNoRows
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&e)
		if strings.TrimSpace(e.Error) != "" {
			return zero, fmt.Errorf("bad status: %s (%s)", resp.Status, e.Error)
		}
		return zero, fmt.Errorf("bad status: %s", resp.Status)
	}

	// для POST/PUT/DELETE обычно неважно что внутри
	_ = json.NewDecoder(resp.Body).Decode(&zero)
	return zero, nil
}

// --------------------
// DTOs
// --------------------

type CategoryDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"`
	CreatedAt string `json:"createdAt"`
}

type SectionDTO struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"`
	CreatedAt          string `json:"createdAt"`
}

type ProductDTO struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`
	OldPrice     *int     `json:"oldPrice"`
	InStock      bool     `json:"inStock"`
	SalePercent  int      `json:"salePercent"`
	ImagePath    string   `json:"imagePath"`
	CreatedAt    string   `json:"createdAt"`
	Badges       []string `json:"badges"`
}

// Requests (CRUD)
type CreateCategoryReq struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"`
}
type UpdateCategoryReq struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"`
}

type CreateSectionReq struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"`
}
type UpdateSectionReq struct {
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"`
}

type CreateProductReq struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`
	InStock      bool     `json:"inStock"`
	ImagePath    string   `json:"imagePath"`
	Badges       []string `json:"badges"`
	Discount     string   `json:"discount"` // sale_20 or ""
}
type UpdateProductReq struct {
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`
	InStock      bool     `json:"inStock"`
	ImagePath    string   `json:"imagePath"`
	Badges       []string `json:"badges"`
	Discount     string   `json:"discount"` // sale_20 or ""
}

// ======================================================
// ✅ GET LISTS: constructions_service:8080 (для pages)
// ======================================================

// func ConstructionsGetCatalogCategories(ctx context.Context, page int, search, orderBy string) ([]CategoryDTO, error) {
// 	u, _ := url.Parse(joinConstructions("/admin/catalog/categories"))
// 	q := u.Query()
// 	if page > 0 {
// 		q.Set("page", fmt.Sprintf("%d", page))
// 	}
// 	if strings.TrimSpace(search) != "" {
// 		q.Set("search", search)
// 	}
// 	if strings.TrimSpace(orderBy) != "" {
// 		q.Set("orderBy", orderBy)
// 	}
// 	u.RawQuery = q.Encode()

// 	return doJSON[[]CategoryDTO](ctx, http.MethodGet, u.String(), nil)
// }

// func ConstructionsGetCatalogSections(ctx context.Context, page int, search, orderBy string) ([]SectionDTO, error) {
// 	u, _ := url.Parse(joinConstructions("/admin/catalog/sections"))
// 	q := u.Query()
// 	if page > 0 {
// 		q.Set("page", fmt.Sprintf("%d", page))
// 	}
// 	if strings.TrimSpace(search) != "" {
// 		q.Set("search", search)
// 	}
// 	if strings.TrimSpace(orderBy) != "" {
// 		q.Set("orderBy", orderBy)
// 	}
// 	u.RawQuery = q.Encode()

// 	return doJSON[[]SectionDTO](ctx, http.MethodGet, u.String(), nil)
// }

// func ConstructionsGetCatalogProducts(ctx context.Context, page int, search, orderBy, categorySlug, sectionSlug string) ([]ProductDTO, error) {
// 	u, _ := url.Parse(joinConstructions("/admin/catalog/products"))
// 	q := u.Query()
// 	if page > 0 {
// 		q.Set("page", fmt.Sprintf("%d", page))
// 	}
// 	if strings.TrimSpace(search) != "" {
// 		q.Set("search", search)
// 	}
// 	if strings.TrimSpace(orderBy) != "" {
// 		q.Set("orderBy", orderBy)
// 	}
// 	if strings.TrimSpace(categorySlug) != "" {
// 		q.Set("categorySlug", categorySlug)
// 	}
// 	if strings.TrimSpace(sectionSlug) != "" {
// 		q.Set("sectionSlug", sectionSlug)
// 	}
// 	u.RawQuery = q.Encode()

// 	return doJSON[[]ProductDTO](ctx, http.MethodGet, u.String(), nil)
// }

// ======================================================
// ✅ CRUD: admin-service (handlers проксируют дальше)
// ======================================================

func AdminCreateCatalogCategory(ctx context.Context, req CreateCategoryReq) error {
	u := joinAdminService("/admin/catalog/categories")
	_, err := doJSON[map[string]any](ctx, http.MethodPost, u, req)
	return err
}
func AdminUpdateCatalogCategory(ctx context.Context, id string, req UpdateCategoryReq) error {
	u := joinAdminService("/admin/catalog/categories/" + url.PathEscape(strings.TrimSpace(id)))
	_, err := doJSON[map[string]any](ctx, http.MethodPut, u, req)
	return err
}
func AdminDeleteCatalogCategory(ctx context.Context, id string) error {
	u := joinAdminService("/admin/catalog/categories/" + url.PathEscape(strings.TrimSpace(id)))
	_, err := doJSON[map[string]any](ctx, http.MethodDelete, u, nil)
	return err
}

func AdminCreateCatalogSection(ctx context.Context, req CreateSectionReq) error {
	u := joinAdminService("/admin/catalog/sections")
	_, err := doJSON[map[string]any](ctx, http.MethodPost, u, req)
	return err
}
func AdminUpdateCatalogSection(ctx context.Context, id string, req UpdateSectionReq) error {
	u := joinAdminService("/admin/catalog/sections/" + url.PathEscape(strings.TrimSpace(id)))
	_, err := doJSON[map[string]any](ctx, http.MethodPut, u, req)
	return err
}
func AdminDeleteCatalogSection(ctx context.Context, id string) error {
	u := joinAdminService("/admin/catalog/sections/" + url.PathEscape(strings.TrimSpace(id)))
	_, err := doJSON[map[string]any](ctx, http.MethodDelete, u, nil)
	return err
}

func AdminCreateCatalogProduct(ctx context.Context, req CreateProductReq) error {
	u := joinAdminService("/admin/catalog/products")
	_, err := doJSON[map[string]any](ctx, http.MethodPost, u, req)
	return err
}
func AdminUpdateCatalogProduct(ctx context.Context, id string, req UpdateProductReq) error {
	u := joinAdminService("/admin/catalog/products/" + url.PathEscape(strings.TrimSpace(id)))
	_, err := doJSON[map[string]any](ctx, http.MethodPut, u, req)
	return err
}
func AdminDeleteCatalogProduct(ctx context.Context, id string) error {
	u := joinAdminService("/admin/catalog/products/" + url.PathEscape(strings.TrimSpace(id)))
	_, err := doJSON[map[string]any](ctx, http.MethodDelete, u, nil)
	return err
}

// Обновите структуры ответов
type PaginatedCategoriesResponse struct {
	Items   []CategoryDTO `json:"items"`
	Total   int           `json:"total"`
	Page    int           `json:"page"`
	PerPage int           `json:"perPage"`
}

type PaginatedSectionsResponse struct {
	Items   []SectionDTO `json:"items"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	PerPage int          `json:"perPage"`
}

type PaginatedProductsResponse struct {
	Items   []ProductDTO `json:"items"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	PerPage int          `json:"perPage"`
}

// Обновите функции
func ConstructionsGetCatalogCategories(ctx context.Context, page int, search, orderBy string) (PaginatedCategoriesResponse, error) {
	u, _ := url.Parse(joinConstructions("/admin/catalog/categories"))
	q := u.Query()
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	q.Set("perPage", "20")
	if strings.TrimSpace(search) != "" {
		q.Set("search", search)
	}
	if strings.TrimSpace(orderBy) != "" {
		q.Set("orderBy", orderBy)
	}
	u.RawQuery = q.Encode()

	return doJSON[PaginatedCategoriesResponse](ctx, http.MethodGet, u.String(), nil)
}

func ConstructionsGetCatalogSections(ctx context.Context, page int, search, orderBy string) (PaginatedSectionsResponse, error) {
	u, _ := url.Parse(joinConstructions("/admin/catalog/sections"))
	q := u.Query()
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	q.Set("perPage", "20")
	if strings.TrimSpace(search) != "" {
		q.Set("search", search)
	}
	if strings.TrimSpace(orderBy) != "" {
		q.Set("orderBy", orderBy)
	}
	u.RawQuery = q.Encode()

	return doJSON[PaginatedSectionsResponse](ctx, http.MethodGet, u.String(), nil)
}

func ConstructionsGetCatalogProducts(ctx context.Context, page int, search, orderBy, categorySlug, sectionSlug string) (PaginatedProductsResponse, error) {
	u, _ := url.Parse(joinConstructions("/admin/catalog/products"))
	q := u.Query()
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	q.Set("perPage", "20")
	if strings.TrimSpace(search) != "" {
		q.Set("search", search)
	}
	if strings.TrimSpace(orderBy) != "" {
		q.Set("orderBy", orderBy)
	}
	if strings.TrimSpace(categorySlug) != "" {
		q.Set("categorySlug", categorySlug)
	}
	if strings.TrimSpace(sectionSlug) != "" {
		q.Set("sectionSlug", sectionSlug)
	}
	u.RawQuery = q.Encode()

	return doJSON[PaginatedProductsResponse](ctx, http.MethodGet, u.String(), nil)
}
