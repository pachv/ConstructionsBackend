package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

// GET /admin/fonts
// Прокси: список шрифтов (если захочешь подгружать через ajax)
func (h *Handler) ProxyAdminFontsList(c *gin.Context) {
	items, err := sender.GetAdminFonts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []sender.AdminFont{}
	}
	c.JSON(http.StatusOK, items)
}

// POST /admin/fonts
// multipart/form-data: name + file
func (h *Handler) ProxyAdminFontsCreate(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// сохраняем во временный файл
	tmpDir := os.TempDir()
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	tmpPath := filepath.Join(tmpDir, "upload_font_*"+ext)

	tmpFile, err := os.CreateTemp(tmpDir, "upload_font_*"+ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp file"})
		return
	}
	tmpPath = tmpFile.Name()
	_ = tmpFile.Close()

	// gin умеет сохранять multipart file на диск
	if err := c.SaveUploadedFile(fh, tmpPath); err != nil {
		_ = os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded file"})
		return
	}
	defer func() { _ = os.Remove(tmpPath) }()

	id, err := sender.CreateAdminFont(c.Request.Context(), name, tmpPath)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// POST /admin/fonts/:id/select
func (h *Handler) ProxyAdminFontsSelect(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := sender.SelectAdminFont(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /admin/fonts/:id
func (h *Handler) ProxyAdminFontsDelete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := sender.DeleteAdminFont(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
