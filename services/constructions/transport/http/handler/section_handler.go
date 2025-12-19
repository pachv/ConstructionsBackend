package handler

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

func (h *Handler) GetSectionsAll(c *gin.Context) {
	items, err := h.siteSectionService.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("GetSectionsAll error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// если пусто — вернуть пустой массив
	if items == nil {
		items = []entity.SiteSectionSummary{}
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) GetSectionBySlug(c *gin.Context) {
	slug := c.Param("slug")

	item, err := h.siteSectionService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		h.logger.Error("GetSectionBySlug error", "err", err, "slug", slug)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GET /api/v1/sections/gallery/picture/:name
func (h *Handler) GetSectionGalleryPicture(c *gin.Context) {
	filename := c.Param("name")
	filePath := filepath.Join("./uploads/sections/gallery", filename)

	fmt.Println("filepath is " + filePath)

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
