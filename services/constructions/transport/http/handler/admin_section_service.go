package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pachv/constructions/constructions/internal/services"
)

func (h *Handler) AdminGetSections(c *gin.Context) {
	page := 1
	if p := strings.TrimSpace(c.Query("page")); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	search := c.Query("search")
	orderBy := c.Query("orderBy")

	items, pageAmount, total, err := h.adminSectionService.GetSectionsSummary(page, search, orderBy)
	if err != nil {
		h.logger.Error("AdminGetSections: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services.AdminSectionsPage{
		Items:      items,
		Page:       page,
		PageAmount: pageAmount,
		Total:      total,
	})
}

func (h *Handler) AdminGetSectionBySlug(c *gin.Context) {
	slug := c.Param("slug")

	section, err := h.adminSectionService.GetSectionBySlug(slug)
	if err != nil {
		msg := err.Error()
		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, section)
}

func (h *Handler) AdminCreateSection(c *gin.Context) {
	var in services.AdminUpsertSectionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.adminSectionService.CreateSection(in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (h *Handler) AdminUpdateSection(c *gin.Context) {
	id := c.Param("id")

	var in services.AdminUpsertSectionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.adminSectionService.UpdateSection(id, in); err != nil {
		msg := err.Error()
		if strings.Contains(strings.ToLower(msg), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
