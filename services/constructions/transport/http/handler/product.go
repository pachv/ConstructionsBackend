package handler

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllCategories(c *gin.Context) {
	items, err := h.productService.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get categories",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetProductPicture(c *gin.Context) {
	filename := c.Param("image")
	filePath := filepath.Join("./uploads/products", filename)

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

func (h *Handler) GetAllSections(c *gin.Context) {
	items, err := h.productService.GetAllSections(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sections"})
		return
	}
	c.JSON(http.StatusOK, items)
}
