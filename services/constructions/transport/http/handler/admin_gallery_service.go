package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pachv/constructions/constructions/internal/services"
)

// ===== api/v1/gallery =====

func (h *Handler) AdminGetGalleryCategories(c *gin.Context) {
	items, err := h.adminGalleryService.ListCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("GetGalleryCategories: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) AdminGetGalleryPhotosByCategory(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	items, err := h.adminGalleryService.ListPhotosByCategorySlug(c.Request.Context(), slug)
	if err != nil {
		h.logger.Error("GetGalleryPhotosByCategory: failed", "err", err, "slug", slug)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list photos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) AdminGetGalleryPicture(c *gin.Context) {
	image := c.Param("image")

	full, err := h.adminGalleryService.OpenPicturePath(image)
	if err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "picture not found"})
			return
		}
		h.logger.Error("GetGalleryPicture: failed", "err", err, "image", image)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.File(full)
}

// ===== admin/gallery =====

func (h *Handler) AdminCreateGalleryCategory(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		// на всякий случай разрешим JSON
		var body struct {
			Title string `json:"title"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			title = strings.TrimSpace(body.Title)
		}
	}
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	cat, err := h.adminGalleryService.CreateCategory(c.Request.Context(), title)
	if err != nil {
		if errors.Is(err, services.ErrGalleryLimitReached) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("AdminCreateGalleryCategory: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"category": cat})
}

func (h *Handler) AdminUpdateGalleryCategory(c *gin.Context) {
	id := c.Param("id")

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		var body struct {
			Title string `json:"title"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			title = strings.TrimSpace(body.Title)
		}
	}
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	cat, err := h.adminGalleryService.UpdateCategory(c.Request.Context(), id, title)
	if err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": cat})
}

func (h *Handler) AdminDeleteGalleryCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.adminGalleryService.DeleteCategory(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		h.logger.Error("AdminDeleteGalleryCategory: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminAddGalleryPhoto(c *gin.Context) {
	categoryID := c.Param("id")

	fh, err := c.FormFile("photo")
	if err != nil || fh == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo is required"})
		return
	}

	alt := c.PostForm("alt")
	sortOrder := c.PostForm("sort_order")

	photo, err := h.adminGalleryService.AddPhoto(c.Request.Context(), categoryID, alt, sortOrder, fh)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGalleryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		case errors.Is(err, services.ErrGalleryBadImage):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			h.logger.Error("AdminAddGalleryPhoto: failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add photo"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"photo": photo})
}

func (h *Handler) AdminDeleteGalleryPhoto(c *gin.Context) {
	id := c.Param("id")

	if err := h.adminGalleryService.DeletePhoto(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
			return
		}
		h.logger.Error("AdminDeleteGalleryPhoto: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
