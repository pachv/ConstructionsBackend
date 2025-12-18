package handler

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetGalleryCategories(c *gin.Context) {
	cats, err := h.galleryService.GetAllCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("GetGalleryCategories failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load gallery categories"})
		return
	}

	// фронту нужны названия (для alt) и slug (чтобы запрашивать фото)
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) GetGalleryPhotosByCategory(c *gin.Context) {
	slug := c.Param("slug")

	photos, err := h.galleryService.GetPhotosByCategorySlug(c.Request.Context(), slug)
	if err != nil {
		h.logger.Error("GetGalleryPhotosByCategory failed", "err", err, "slug", slug)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load gallery photos"})
		return
	}

	// отдаём alt + ссылку на изображение (image_path)
	c.JSON(http.StatusOK, photos)
}

func (h *Handler) GetGalleryPicture(c *gin.Context) {
	filename := c.Param("image")
	filePath := filepath.Join("./uploads/gallery", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+filename)
	c.File(filePath)
}
