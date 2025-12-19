package handler

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

//
// =========================
// CONSTANTS
// =========================
//

const (
	SectionsMainImagesDir    = "./img/sections/main"
	SectionsGalleryImagesDir = "./img/sections/gallery"
	CatalogImagesDir         = "./img/catalog"
)

//
// =========================
// INTERNAL HELPERS
// =========================
//

func sanitizePath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = filepath.Clean(p)
	p = strings.TrimPrefix(p, "/")

	if p == "" || strings.HasPrefix(p, "..") {
		return ""
	}
	return p
}

func serveImage(c *gin.Context, baseDir string, rawPath string) {
	path := sanitizePath(rawPath)
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad filename"})
		return
	}

	fullPath := filepath.Join(baseDir, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(fullPath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+filepath.Base(fullPath))
	c.File(fullPath)
}

//
// =========================
// HANDLERS (METHODS OF *Handler)
// =========================
//

// GET /api/v1/sections/picture/:name
func (h *Handler) GetSectionMainPicture(c *gin.Context) {
	serveImage(c, SectionsMainImagesDir, c.Param("name"))
}

// GET /api/v1/sections/gallery/picture/:name
func (h *Handler) GetSectionGalleryPicture(c *gin.Context) {
	serveImage(c, SectionsGalleryImagesDir, c.Param("name"))
}

// GET /api/v1/catalog/picture/:name
func (h *Handler) GetCatalogPicture(c *gin.Context) {
	serveImage(c, CatalogImagesDir, c.Param("name"))
}

//
// =========================
// GET /api/v1/sections
// =========================
//

func (h *Handler) GetSectionsAll(c *gin.Context) {
	items, err := h.siteSectionService.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("GetSectionsAll error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if items == nil {
		items = []entity.SiteSectionSummary{}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

//
// =========================
// GET /api/v1/sections/:slug
// =========================
//

func (h *Handler) GetSectionBySlug(c *gin.Context) {
	slug := c.Param("slug")

	section, err := h.siteSectionService.GetBySlugFull(c.Request.Context(), slug)
	if err != nil {
		h.logger.Error("GetSectionBySlug error", "slug", slug, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if section == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "section not found"})
		return
	}

	// гарантируем поля, чтобы фронт не ебался
	if section.AdvantegesText == "" {
		section.AdvantegesText = ""
	}
	if section.AdvantegesArray == nil {
		section.AdvantegesArray = []string{}
	}
	if section.Gallery == nil {
		section.Gallery = []entity.SiteSectionGallery{}
	}
	if section.HasCatalog && section.Catalog != nil {
		if section.Catalog.Categories == nil {
			section.Catalog.Categories = []entity.SiteSectionCatalogCategory{}
		}
		if section.Catalog.Items == nil {
			section.Catalog.Items = []entity.SiteSectionCatalogItem{}
		}
	}

	c.JSON(http.StatusOK, section)
}
