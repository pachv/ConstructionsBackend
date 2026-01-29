package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// constructions service (куда реально уходит CRUD)
const constructionsCatalogBase = "http://constructions_service:8080/admin/catalog"

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: 25 * time.Second}
}

// ===============
// helpers
// ===============

func boolFromForm(v string) bool {
	v = strings.TrimSpace(v)
	return v == "true" || v == "1" || strings.EqualFold(v, "on") || strings.EqualFold(v, "yes")
}

func intFromForm(v string) (int, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}

func (h *Handler) proxyMultipart(
	c *gin.Context,
	method string,
	targetURL string,
	fields map[string]string,
	outFileField string, // как называется поле файла на constructions (например imagePath)
	incomingFileField string, // как называется input file у браузера (например image)
) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	for k, v := range fields {
		_ = w.WriteField(k, v)
	}

	if outFileField != "" && incomingFileField != "" {
		fh, err := c.FormFile(incomingFileField)
		if err == nil && fh != nil {
			f, err := fh.Open()
			if err != nil {
				_ = w.Close()
				c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
				return
			}
			defer f.Close()

			part, err := w.CreateFormFile(outFileField, filepath.Base(fh.Filename))
			if err != nil {
				_ = w.Close()
				c.JSON(http.StatusBadRequest, gin.H{"error": "create form file: " + err.Error()})
				return
			}
			_, _ = io.Copy(part, f)
		}
	}

	_ = w.Close()

	req, err := http.NewRequest(method, targetURL, &body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "new request: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := newHTTPClient().Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "proxy request failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		// возвращаем как есть
		c.Data(resp.StatusCode, "application/json", b)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) proxyJSONN(c *gin.Context, method, targetURL string, payload any) {
	var rdr io.Reader
	if payload != nil {
		bin, _ := json.Marshal(payload)
		rdr = bytes.NewReader(bin)
	}

	req, err := http.NewRequest(method, targetURL, rdr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "new request: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := newHTTPClient().Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "proxy request failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		c.Data(resp.StatusCode, "application/json", b)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// =========================
// SECTIONS: from page -> constructions
// endpoints in admin-service:
// POST /admin-service/admin/products/sections/create
// POST /admin-service/admin/products/sections/update
// POST /admin-service/admin/products/sections/delete
// =========================

func (h *Handler) CreateCatalogSectionFromPage(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	parent := strings.TrimSpace(c.PostForm("parentCategorySlug"))

	if title == "" || parent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and parentCategorySlug required"})
		return
	}

	target := constructionsCatalogBase + "/sections"
	h.proxyMultipart(
		c,
		http.MethodPost,
		target,
		map[string]string{
			"title":              title,
			"parentCategorySlug": parent,
			// slug/id пусть генерирует constructions-service (если у тебя так сделано)
		},
		"imagePath",
		"image",
	)
}

func (h *Handler) UpdateCatalogSectionFromPage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	title := strings.TrimSpace(c.PostForm("title"))
	parent := strings.TrimSpace(c.PostForm("parentCategorySlug"))

	if id == "" || title == "" || parent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id, title, parentCategorySlug required"})
		return
	}

	target := constructionsCatalogBase + "/sections/" + url.PathEscape(id)
	h.proxyMultipart(
		c,
		http.MethodPut,
		target,
		map[string]string{
			"title":              title,
			"parentCategorySlug": parent,
		},
		"imagePath",
		"image", // может быть пусто — это ок, файл просто не будет отправлен
	)
}

func (h *Handler) DeleteCatalogSectionFromPage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	target := constructionsCatalogBase + "/sections/" + url.PathEscape(id)
	h.proxyJSON(c, http.MethodDelete, target, nil)
}

// =========================
// CATEGORIES: from page -> constructions
// endpoints in admin-service:
// POST /admin-service/admin/products/categories/create
// POST /admin-service/admin/products/categories/update
// POST /admin-service/admin/products/categories/delete
// =========================

func (h *Handler) CreateCatalogCategoryFromPage(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}

	target := constructionsCatalogBase + "/categories"
	h.proxyMultipart(
		c,
		http.MethodPost,
		target,
		map[string]string{
			"title": title,
		},
		"imagePath",
		"image",
	)
}

func (h *Handler) UpdateCatalogCategoryFromPage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	title := strings.TrimSpace(c.PostForm("title"))

	if id == "" || title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and title required"})
		return
	}

	target := constructionsCatalogBase + "/categories/" + url.PathEscape(id)
	h.proxyMultipart(
		c,
		http.MethodPut,
		target,
		map[string]string{
			"title": title,
		},
		"imagePath",
		"image",
	)
}

func (h *Handler) DeleteCatalogCategoryFromPage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	target := constructionsCatalogBase + "/categories/" + url.PathEscape(id)
	h.proxyJSON(c, http.MethodDelete, target, nil)
}

// =========================
// PRODUCTS: from page -> constructions
// endpoints in admin-service:
// POST /admin-service/admin/products/products/create
// POST /admin-service/admin/products/products/update
// POST /admin-service/admin/products/products/delete
//
// ВАЖНО: discount ожидаем как "sale_20" (или пусто).
// badges[] пробрасываем массивом.
// =========================

func (h *Handler) CreateCatalogProductFromPage(c *gin.Context) {

	title := strings.TrimSpace(c.PostForm("title"))
	categorySlug := strings.TrimSpace(c.PostForm("category"))
	sectionSlug := strings.TrimSpace(c.PostForm("section"))
	typ := strings.TrimSpace(c.PostForm("type"))
	brand := strings.TrimSpace(c.PostForm("brand"))
	discount := strings.TrimSpace(c.PostForm("discount"))

	priceStr := strings.TrimSpace(c.PostForm("price"))

	if title == "" || categorySlug == "" || sectionSlug == "" || typ == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, categorySlug, sectionSlug, type, price required"})
		return
	}

	// multipart: обычные поля + badges[] + файл
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	_ = w.WriteField("title", title)
	_ = w.WriteField("categorySlug", categorySlug)
	_ = w.WriteField("sectionSlug", sectionSlug)
	_ = w.WriteField("type", typ)
	_ = w.WriteField("brand", brand)
	_ = w.WriteField("price", priceStr)
	_ = w.WriteField("inStock", strconv.FormatBool(boolFromForm(c.PostForm("inStock"))))
	_ = w.WriteField("discount", discount)

	for _, b := range c.PostFormArray("badges[]") {
		b = strings.TrimSpace(b)
		if b != "" {
			_ = w.WriteField("badges[]", b)
		}
	}

	fh, err := c.FormFile("image")
	if err == nil && fh != nil {
		f, err := fh.Open()
		if err != nil {
			_ = w.Close()
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		part, err := w.CreateFormFile("imagePath", filepath.Base(fh.Filename))
		if err != nil {
			_ = w.Close()
			c.JSON(http.StatusBadRequest, gin.H{"error": "create form file: " + err.Error()})
			return
		}
		_, _ = io.Copy(part, f)
	}

	_ = w.Close()

	target := constructionsCatalogBase + "/products"
	req, err := http.NewRequest(http.MethodPost, target, &body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "new request: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := newHTTPClient().Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "proxy request failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	bin, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		c.Data(resp.StatusCode, "application/json", bin)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) UpdateCatalogProductFromPage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	price, err := intFromForm(c.PostForm("price"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad price"})
		return
	}

	payload := map[string]any{
		"title":        strings.TrimSpace(c.PostForm("title")),
		"slug":         strings.TrimSpace(c.PostForm("slug")), // можно пусто, если на сервисе генеришь по title
		"categorySlug": strings.TrimSpace(c.PostForm("categorySlug")),
		"sectionSlug":  strings.TrimSpace(c.PostForm("sectionSlug")),
		"brand":        strings.TrimSpace(c.PostForm("brand")),
		"type":         strings.TrimSpace(c.PostForm("type")),
		"price":        price,
		"inStock":      boolFromForm(c.PostForm("inStock")),
		"discount":     strings.TrimSpace(c.PostForm("discount")), // "" => revert
		"badges":       c.PostFormArray("badges[]"),
		"imagePath":    "", // update image отдельно, если надо (у тебя service сейчас без multipart update products)
	}

	target := constructionsCatalogBase + "/products/" + url.PathEscape(id)
	h.proxyJSONN(c, http.MethodPut, target, payload)
}

func (h *Handler) DeleteCatalogProductFromPage(c *gin.Context) {
	id := strings.TrimSpace(c.PostForm("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	target := constructionsCatalogBase + "/products/" + url.PathEscape(id)
	h.proxyJSON(c, http.MethodDelete, target, nil)
}

// маленький sanity-check чтобы быстрее ловить пустые обязательные поля при update
func ensureNotEmpty(name, v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("%s required", name)
	}
	return nil
}
