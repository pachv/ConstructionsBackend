package handler

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllCertificates(c *gin.Context) {
	items, err := h.certificateService.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("GetAllCertificates: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get certificates"})
		return
	}

	// отдаём как есть: title + file_path
	c.JSON(http.StatusOK, items)
}

// GET /certificates/file/:name
func (h *Handler) GetCertificateFile(c *gin.Context) {
	filename := c.Param("name")
	filePath := filepath.Join("./uploads/certificates", filename)

	h.logger.Debug("GetCertificateFile: request", "filename", filename, "path", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		h.logger.Error("GetCertificateFile: stat failed", "err", err, "path", filePath)
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
