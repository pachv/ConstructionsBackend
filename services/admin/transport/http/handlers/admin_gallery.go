package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const constructionsGalleryAdminBaseURL = "http://constructions_service:8080/admin/gallery"

// POST /admin-service/admin/gallery/categories
// body: {"title":"..."}
func (h *Handler) GalleryCreateCategoryProxy(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)

	req, err := http.NewRequest(http.MethodPost, constructionsGalleryAdminBaseURL+"/categories", bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", out)
}

// PUT /admin-service/admin/gallery/categories/:id
// body: {"title":"..."}
func (h *Handler) GalleryUpdateCategoryProxy(c *gin.Context) {
	id := c.Param("id")
	body, _ := io.ReadAll(c.Request.Body)

	req, err := http.NewRequest(http.MethodPut, constructionsGalleryAdminBaseURL+"/categories/"+id, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", out)
}

// DELETE /admin-service/admin/gallery/categories/:id
func (h *Handler) GalleryDeleteCategoryProxy(c *gin.Context) {
	id := c.Param("id")

	req, err := http.NewRequest(http.MethodDelete, constructionsGalleryAdminBaseURL+"/categories/"+id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", out)
}

// POST /admin-service/admin/gallery/categories/:id/photos
// multipart/form-data: photo=<file>  (alt/sort_order можно будет добавить позже)
func (h *Handler) GalleryAddPhotoProxy(c *gin.Context) {
	id := c.Param("id")

	file, err := c.FormFile("photo")
	if err != nil || file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer src.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("photo", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form file"})
		return
	}

	if _, err := io.Copy(part, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to copy file"})
		return
	}

	// если позже захочешь alt/sort_order — раскомментируешь:
	// _ = writer.WriteField("alt", c.PostForm("alt"))
	// _ = writer.WriteField("sort_order", c.PostForm("sort_order"))

	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, constructionsGalleryAdminBaseURL+"/categories/"+id+"/photos", &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", out)
}

// DELETE /admin-service/admin/gallery/photos/:id
func (h *Handler) GalleryDeletePhotoProxy(c *gin.Context) {
	id := c.Param("id")

	req, err := http.NewRequest(http.MethodDelete, constructionsGalleryAdminBaseURL+"/photos/"+id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact constructions service"})
		return
	}
	defer resp.Body.Close()

	out, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", out)
}
