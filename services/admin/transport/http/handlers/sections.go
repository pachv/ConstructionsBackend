package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// =========================
// SECTIONS - LIST & GET
// =========================

// список (GET /inside/sections/all)
func (h *Handler) ProxySectionsList(c *gin.Context) {
	h.proxyyJSON(c, http.MethodGet, constructionsBaseURL+"/admin/sections/all", nil)
}

// секция FULL (GET /inside/sections/view/:slug)
func (h *Handler) ProxySectionBySlug(c *gin.Context) {
	slug := c.Param("slug")
	h.proxyyJSON(c, http.MethodGet, constructionsBaseURL+"/admin/sections/"+slug+"/full", nil)
}

// =========================
// SECTIONS - CREATE
// =========================

// CREATE basic (POST /inside/sections/create)
func (h *Handler) CreateSectionBasicProxy(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image is required"})
		return
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("title", c.PostForm("title"))
	_ = w.WriteField("advantegesText", c.PostForm("advantegesText"))

	advs := c.PostFormArray("advanteges[]")
	if len(advs) == 0 {
		advs = c.PostFormArray("advanteges")
	}
	for _, a := range advs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		_ = w.WriteField("advanteges[]", a)
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	part, err := w.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(500, gin.H{"error": "failed to copy"})
		return
	}

	_ = w.Close()

	req, err := http.NewRequest(http.MethodPost, constructionsBaseURL+"/admin/sections/create-form", &buf)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// =========================
// SECTIONS - UPDATE
// =========================

// UPDATE section (PUT /inside/sections/update/:id)
func (h *Handler) UpdateSectionProxy(c *gin.Context) {
	id := c.Param("id")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("title", c.PostForm("title"))
	_ = w.WriteField("slug", c.PostForm("slug"))
	_ = w.WriteField("advantegesText", c.PostForm("advantegesText"))
	_ = w.WriteField("advantegesArray", c.PostForm("advantegesArray"))
	_ = w.WriteField("hasGallery", c.PostForm("hasGallery"))
	_ = w.WriteField("hasCatalog", c.PostForm("hasCatalog"))

	// если есть файл
	file, err := c.FormFile("image")
	if err == nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()
			part, err := w.CreateFormFile("image", file.Filename)
			if err == nil {
				_, _ = io.Copy(part, src)
			}
		}
	}

	_ = w.Close()

	req, err := http.NewRequest(http.MethodPut, constructionsBaseURL+"/admin/sections/"+id, &buf)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// =========================
// SECTIONS - DELETE
// =========================

// DELETE секции (DELETE /inside/sections/delete/:id)
func (h *Handler) DeleteSectionProxy(c *gin.Context) {
	id := c.Param("id")
	h.proxyyJSON(c, http.MethodDelete, constructionsBaseURL+"/admin/sections/"+id, nil)
}

// =========================
// HELPER: slug -> sectionID
// =========================

func (h *Handler) getSectionIDBySlug(slug string) (string, error) {
	resp, err := http.Get(constructionsBaseURL + "/admin/sections/" + slug + "/full")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", http.ErrAbortHandler
	}

	var section struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&section); err != nil {
		return "", err
	}

	return section.ID, nil
}

// =========================
// GALLERY
// =========================

// ADD gallery photo (POST /inside/sections/gallery/:slug/add)
func (h *Handler) AddGalleryPhotoProxy(c *gin.Context) {
	slug := c.Param("slug")

	// 1. Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	// 2. Отправляем запрос с sectionID
	h.proxyyJSON(c, http.MethodPost, constructionsBaseURL+"/admin/sections/"+sectionID+"/gallery", payload)
}

// DELETE gallery photo (DELETE /inside/sections/gallery/:slug/photo/:photoId)
func (h *Handler) DeleteGalleryPhotoProxy(c *gin.Context) {
	slug := c.Param("slug")
	photoID := c.Param("photoId")

	// Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	h.proxyyJSON(c, http.MethodDelete, constructionsBaseURL+"/admin/sections/"+sectionID+"/gallery/"+photoID, nil)
}

// UPLOAD gallery image (POST /inside/sections/gallery/:slug/upload)
func (h *Handler) UploadGalleryImageProxy(c *gin.Context) {
	slug := c.Param("slug")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image is required"})
		return
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	part, err := w.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(500, gin.H{"error": "failed to copy"})
		return
	}

	_ = w.Close()

	// Отправляем на эндпоинт с slug (это upload, не требует ID)
	req, err := http.NewRequest(http.MethodPost, constructionsBaseURL+"/admin/sections/gallery/"+slug+"/upload", &buf)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// =========================
// CATALOG - CATEGORIES
// =========================

// ADD catalog category (POST /inside/sections/catalog/:slug/categories/add)
func (h *Handler) AddCatalogCategoryProxy(c *gin.Context) {
	slug := c.Param("slug")

	// Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	h.proxyyJSON(c, http.MethodPost, constructionsBaseURL+"/admin/sections/"+sectionID+"/catalog/categories", payload)
}

// DELETE catalog category (DELETE /inside/sections/catalog/:slug/categories/:categoryId)
func (h *Handler) DeleteCatalogCategoryProxy(c *gin.Context) {
	slug := c.Param("slug")
	categoryID := c.Param("categoryId")

	// Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	h.proxyyJSON(c, http.MethodDelete, constructionsBaseURL+"/admin/sections/"+sectionID+"/catalog/categories/"+categoryID, nil)
}

// UPDATE catalog (PUT /inside/sections/catalog/:slug/update)
func (h *Handler) UpdateCatalogProxy(c *gin.Context) {
	slug := c.Param("slug")

	// Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	h.proxyyJSON(c, http.MethodPut, constructionsBaseURL+"/admin/sections/"+sectionID+"/catalog", payload)
}

// =========================
// CATALOG - ITEMS
// =========================

// ADD catalog item (POST /inside/sections/catalog/:slug/items/add)
func (h *Handler) AddCatalogItemProxy(c *gin.Context) {
	slug := c.Param("slug")

	// Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	h.proxyyJSON(c, http.MethodPost, constructionsBaseURL+"/admin/sections/"+sectionID+"/catalog/items", payload)
}

// DELETE catalog item (DELETE /inside/sections/catalog/:slug/items/:itemId)
func (h *Handler) DeleteCatalogItemProxy(c *gin.Context) {
	slug := c.Param("slug")
	itemID := c.Param("itemId")

	// Получаем sectionID по slug
	sectionID, err := h.getSectionIDBySlug(slug)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to get section by slug"})
		return
	}

	h.proxyyJSON(c, http.MethodDelete, constructionsBaseURL+"/admin/sections/"+sectionID+"/catalog/items/"+itemID, nil)
}

// UPLOAD catalog item image (POST /inside/sections/catalog/:slug/items/upload)
func (h *Handler) UploadCatalogItemImageProxy(c *gin.Context) {
	slug := c.Param("slug")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image is required"})
		return
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	part, err := w.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(500, gin.H{"error": "failed to copy"})
		return
	}

	_ = w.Close()

	// Upload endpoint с slug
	req, err := http.NewRequest(http.MethodPost, constructionsBaseURL+"/admin/sections/"+slug+"/catalog/items/upload", &buf)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// =========================
// HELPER
// =========================

func (h *Handler) proxyyJSON(c *gin.Context, method, url string, payload interface{}) {
	var reqBody io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to marshal payload"})
			return
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// TOGGLE Gallery (PATCH /inside/sections/:id/gallery/toggle)
func (h *Handler) ToggleSectionGalleryProxy(c *gin.Context) {
	id := c.Param("id")

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	h.proxyyJSON(c, http.MethodPatch, constructionsBaseURL+"/admin/sections/"+id+"/gallery/toggle", payload)
}

// TOGGLE Catalog (PATCH /inside/sections/:id/catalog/toggle)
func (h *Handler) ToggleSectionCatalogProxy(c *gin.Context) {
	id := c.Param("id")

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	h.proxyyJSON(c, http.MethodPatch, constructionsBaseURL+"/admin/sections/"+id+"/catalog/toggle", payload)
}
