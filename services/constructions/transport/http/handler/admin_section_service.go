package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/services"
)

//
// =========================
// TOGGLES
// =========================
//

// POST /admin/sections/create-form
// multipart: title, advantegesText, advanteges[], image (file)
func (h *Handler) AdminCreateSectionForm(c *gin.Context) {
	// 1. Парсим multipart
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	advText := strings.TrimSpace(c.PostForm("advantegesText"))

	// advanteges[] - массив строк
	advs := c.Request.PostForm["advanteges[]"]
	if advs == nil {
		advs = []string{}
	}

	// 2. Получаем файл
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	savePath := filepath.Join("./uploads/sections/main", filename) // adjust path as needed

	// Создаём директорию если не существует
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		h.logger.Error("AdminCreateSectionForm: failed to create dir", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// Сохраняем файл
	out, err := os.Create(savePath)
	if err != nil {
		h.logger.Error("AdminCreateSectionForm: failed to create file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		h.logger.Error("AdminCreateSectionForm: failed to copy file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// 4. Вызываем сервис
	input := services.AdminCreateSectionFormInput{
		Title:          title,
		ImageFilename:  filename, // только имя файла, не полный путь
		AdvantegesText: advText,
		Advanteges:     advs,
	}

	sectionID, err := h.adminSectionService.CreateSectionFromForm(c.Request.Context(), input)
	if err != nil {
		msg := err.Error()
		h.logger.Error("AdminCreateSectionForm: failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ok": true,
		"id": sectionID,
	})
}

// DELETE /admin/sections/:id
func (h *Handler) AdminDeleteSection(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.adminSectionService.DeleteSection(c.Request.Context(), sectionID); err != nil {
		msg := err.Error()
		h.logger.Error("AdminDeleteSection: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// PATCH /admin/sections/:id/gallery/toggle
// body: { "enabled": true }
func (h *Handler) AdminToggleSectionGallery(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var in services.AdminToggleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.adminSectionService.ToggleGallery(c.Request.Context(), sectionID, in.Enabled); err != nil {
		msg := err.Error()
		h.logger.Error("AdminToggleSectionGallery: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// PATCH /admin/sections/:id/catalog/toggle
// body: { "enabled": true }
func (h *Handler) AdminToggleSectionCatalog(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var in services.AdminToggleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.adminSectionService.ToggleCatalog(c.Request.Context(), sectionID, in.Enabled); err != nil {
		msg := err.Error()
		h.logger.Error("AdminToggleSectionCatalog: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

//
// =========================
// GALLERY
// =========================
//

// POST /admin/sections/:id/gallery
// body: { "name": "...", "url": "baths-1.jpg", "sortOrder": 1 }
func (h *Handler) AdminAddSectionGalleryPhoto(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var in services.AdminAddGalleryPhotoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	photoID, err := h.adminSectionService.AddGalleryPhoto(c.Request.Context(), sectionID, in)
	if err != nil {
		msg := err.Error()
		h.logger.Error("AdminAddSectionGalleryPhoto: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ok": true,
		"id": photoID,
	})
}

// DELETE /admin/sections/:id/gallery/:photoId
func (h *Handler) AdminDeleteSectionGalleryPhoto(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	photoID := strings.TrimSpace(c.Param("photoId"))

	if sectionID == "" || photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and photoId are required"})
		return
	}

	if err := h.adminSectionService.DeleteGalleryPhoto(c.Request.Context(), sectionID, photoID); err != nil {
		msg := err.Error()
		h.logger.Error("AdminDeleteSectionGalleryPhoto: failed", "sectionID", sectionID, "photoID", photoID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

//
// =========================
// CATALOG: CATEGORIES
// =========================
//

// POST /admin/sections/:id/catalog/categories
// body: { "categoryId": "...", "sortOrder": 1 }
func (h *Handler) AdminAddSectionCatalogCategory(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var in services.AdminAddCatalogCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	id, err := h.adminSectionService.AddCatalogCategory(c.Request.Context(), sectionID, in)
	if err != nil {
		msg := err.Error()
		h.logger.Error("AdminAddSectionCatalogCategory: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ok": true,
		"id": id,
	})
}

// DELETE /admin/sections/:id/catalog/categories/:categoryId
func (h *Handler) AdminDeleteSectionCatalogCategory(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	categoryID := strings.TrimSpace(c.Param("categoryId"))

	if sectionID == "" || categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and categoryId are required"})
		return
	}

	if err := h.adminSectionService.DeleteCatalogCategory(c.Request.Context(), sectionID, categoryID); err != nil {
		msg := err.Error()
		h.logger.Error("AdminDeleteSectionCatalogCategory: failed", "sectionID", sectionID, "categoryID", categoryID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

//
// =========================
// CATALOG: ITEMS
// =========================
//

// POST /admin/sections/:id/catalog/items
func (h *Handler) AdminAddSectionCatalogItem(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var in services.AdminAddCatalogItemInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	itemID, err := h.adminSectionService.AddCatalogItem(c.Request.Context(), sectionID, in)
	if err != nil {
		msg := err.Error()
		h.logger.Error("AdminAddSectionCatalogItem: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ok": true,
		"id": itemID,
	})
}

// DELETE /admin/sections/:id/catalog/items/:itemId
func (h *Handler) AdminDeleteSectionCatalogItem(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	itemID := strings.TrimSpace(c.Param("itemId"))

	if sectionID == "" || itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id and itemId are required"})
		return
	}

	if err := h.adminSectionService.DeleteCatalogItem(c.Request.Context(), sectionID, itemID); err != nil {
		msg := err.Error()
		h.logger.Error("AdminDeleteSectionCatalogItem: failed", "sectionID", sectionID, "itemID", itemID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type adminSectionsListResponse struct {
	Items []adminSectionSummaryDTO `json:"items"`
}

type adminSectionSummaryDTO struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Label      string `json:"label"`
	Slug       string `json:"slug"`
	Image      string `json:"image"`
	HasGallery bool   `json:"hasGallery"`
	HasCatalog bool   `json:"hasCatalog"`
}

// GET /admin/sections/all
func (h *Handler) AdminGetSectionsAll(c *gin.Context) {
	items, err := h.adminSectionService.GetAllSectionsSummary()
	if err != nil {
		h.logger.Error("AdminGetSectionsAll: failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out := adminSectionsListResponse{Items: make([]adminSectionSummaryDTO, 0, len(items))}
	for _, it := range items {
		if it == nil {
			continue
		}
		out.Items = append(out.Items, adminSectionSummaryDTO{
			ID:         it.ID,
			Title:      it.Title,
			Label:      it.Label,
			Slug:       it.Slug,
			Image:      it.Image,
			HasGallery: it.HasGallery,
			HasCatalog: it.HasCatalog,
		})
	}

	c.JSON(http.StatusOK, out)
}

// GET /admin/sections/:slug/full
func (h *Handler) AdminGetSectionFullBySlug(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	sec, err := h.adminSectionService.GetSectionFullBySlug(slug)
	if err != nil {
		msg := err.Error()
		h.logger.Error("AdminGetSectionFullBySlug: failed", "slug", slug, "err", err)

		// у тебя сервис возвращает "section not found"
		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// ⚠️ Важно: твой service уже формирует URL картинок и заполняет структуры как нужно.
	// Просто отдаем объект секции.
	c.JSON(http.StatusOK, sec)
}

func (h *Handler) AdminUpdateSectionForm(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 1) multipart
	if err := c.Request.ParseMultipartForm(20 << 20); err != nil { // 20MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	slug := strings.TrimSpace(c.PostForm("slug"))
	advText := strings.TrimSpace(c.PostForm("advantegesText"))

	// advantegesArray приходит как JSON строка (как у тебя на фронте)
	rawAdvs := strings.TrimSpace(c.PostForm("advantegesArray"))
	advs := make([]string, 0)

	if rawAdvs != "" {
		if err := json.Unmarshal([]byte(rawAdvs), &advs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "advantegesArray must be valid json array"})
			return
		}
		// чистим пустые
		clean := make([]string, 0, len(advs))
		for _, a := range advs {
			a = strings.TrimSpace(a)
			if a != "" {
				clean = append(clean, a)
			}
		}
		advs = clean
	}

	// 2) optional image file
	var filename string // пусто = не меняем картинку
	file, header, err := c.Request.FormFile("image")
	if err == nil && file != nil && header != nil {
		defer file.Close()

		filename = fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
		savePath := filepath.Join("./uploads/sections/main", filename)

		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			h.logger.Error("AdminUpdateSectionForm: failed to create dir", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}

		out, err := os.Create(savePath)
		if err != nil {
			h.logger.Error("AdminUpdateSectionForm: failed to create file", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			h.logger.Error("AdminUpdateSectionForm: failed to copy file", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}
	}

	// 3) сервис
	in := services.AdminUpdateSectionFormInput{
		ID:             sectionID,
		Title:          title,
		Slug:           slug,
		AdvantegesText: advText,
		Advanteges:     advs,
		ImageFilename:  filename, // "" => не менять
	}

	if err := h.adminSectionService.UpdateSectionFromForm(c.Request.Context(), in); err != nil {
		msg := err.Error()
		h.logger.Error("AdminUpdateSectionForm: failed", "sectionID", sectionID, "err", err)

		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /admin/sections/:slug/gallery/upload
// Загружает картинку и возвращает URL/имя файла
func (h *Handler) AdminUploadGalleryImage(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	// Генерируем имя файла
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	savePath := filepath.Join("./uploads/sections/gallery", filename)

	// Создаём директорию
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		h.logger.Error("AdminUploadGalleryImage: failed to create dir", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// Сохраняем файл
	out, err := os.Create(savePath)
	if err != nil {
		h.logger.Error("AdminUploadGalleryImage: failed to create file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		h.logger.Error("AdminUploadGalleryImage: failed to copy file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// Возвращаем только имя файла (или относительный URL)
	c.JSON(http.StatusOK, gin.H{
		"url": filename,
	})
}

// POST /admin/sections/:slug/catalog/items/upload
// Загружает картинку товара и возвращает URL/имя файла
func (h *Handler) AdminUploadCatalogItemImage(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	// Генерируем имя файла
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	savePath := filepath.Join("./uploads/catalog", filename)

	// Создаём директорию
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		h.logger.Error("AdminUploadCatalogItemImage: failed to create dir", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// Сохраняем файл
	out, err := os.Create(savePath)
	if err != nil {
		h.logger.Error("AdminUploadCatalogItemImage: failed to create file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		h.logger.Error("AdminUploadCatalogItemImage: failed to copy file", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// Возвращаем только имя файла
	c.JSON(http.StatusOK, gin.H{
		"url": filename,
	})
}
