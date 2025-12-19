package handler

import (
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminGetAllCertificates(c *gin.Context) {
	items, err := h.certificateAdminService.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminGetAllCertificates: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get certificates"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) AdminCreateCertificate(c *gin.Context) {
	title := c.PostForm("title")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	created, err := h.certificateAdminService.Create(c.Request.Context(), title, file)
	if err != nil {
		h.logger.Error("AdminCreateCertificate: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create certificate"})
		return
	}

	c.JSON(http.StatusOK, created)
}

func (h *Handler) AdminUpdateCertificate(c *gin.Context) {
	id := c.Param("id")

	var titlePtr *string
	if t, ok := c.GetPostForm("title"); ok {
		tt := t
		titlePtr = &tt
	}

	// file optional
	var filePtr *multipart.FileHeader
	if f, err := c.FormFile("file"); err == nil && f != nil {
		filePtr = f
	}

	updated, err := h.certificateAdminService.Update(c.Request.Context(), id, titlePtr, filePtr)
	if err != nil {
		h.logger.Error("AdminUpdateCertificate: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update certificate"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *Handler) AdminDeleteCertificate(c *gin.Context) {
	id := c.Param("id")

	if err := h.certificateAdminService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("AdminDeleteCertificate: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /admin/certificates/file/:name
func (h *Handler) AdminGetCertificateFile(c *gin.Context) {
	filename := c.Param("name")
	filePath := filepath.Join("./uploads/certificates", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		h.logger.Error("AdminGetCertificateFile: stat failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot access file"})
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	c.File(filePath)
}
