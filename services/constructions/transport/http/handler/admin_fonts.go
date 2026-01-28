package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type adminCreateFontResp struct {
	ID string `json:"id"`
}

type adminCreateFontForm struct {
	Name string `form:"name"`
}

// GET /admin/fonts
func (h *Handler) AdminListFonts(c *gin.Context) {
	items, err := h.adminFontsService.ListFonts(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminListFonts", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list fonts"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /admin/fonts (multipart/form-data: name + file)
func (h *Handler) AdminCreateFont(c *gin.Context) {
	var form adminCreateFontForm
	_ = c.ShouldBind(&form)

	form.Name = strings.TrimSpace(form.Name)
	if form.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	row, err := h.adminFontsService.CreateFontFromForm(c.Request.Context(), form.Name, fh)
	if err != nil {
		h.logger.Error("AdminCreateFont", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adminCreateFontResp{ID: row.ID})
}

// DELETE /admin/fonts/:id
func (h *Handler) AdminDeleteFont(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.adminFontsService.DeleteFont(c.Request.Context(), id); err != nil {
		h.logger.Error("AdminDeleteFont", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete font"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /admin/fonts/:id/select
func (h *Handler) AdminSelectFont(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.adminFontsService.SelectFont(c.Request.Context(), id); err != nil {
		h.logger.Error("AdminSelectFont", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to select font"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
