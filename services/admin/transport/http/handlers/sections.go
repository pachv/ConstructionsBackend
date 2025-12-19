package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const PublicAPIBaseURL = "http://localhost:80"

func (h *Handler) ProxySectionsList(c *gin.Context) {
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURL+"/api/v1/sections", nil)
}

func (h *Handler) ProxySectionBySlug(c *gin.Context) {
	slug := c.Param("slug")
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURL+"/api/v1/sections/"+slug, nil)
}

func (h *Handler) AddGalleryItem(c *gin.Context) {
	slug := c.Param("slug")
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "cant read body"})
		return
	}
	h.proxyJSON(c, http.MethodPost, PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/gallery", raw)
}

func (h *Handler) DeleteGalleryItem(c *gin.Context) {
	slug := c.Param("slug")
	id := c.Param("id")
	h.proxyJSON(c, http.MethodDelete, PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/gallery/"+id, nil)
}

func (h *Handler) UploadGalleryPicture(c *gin.Context) {
	slug := c.Param("slug")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "no file provided"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// можно передать имя файла, slug и т.п.
	part, err := w.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(500, gin.H{"error": "failed to copy"})
		return
	}

	_ = w.WriteField("slug", slug)
	w.Close()

	req, err := http.NewRequest(http.MethodPost,
		PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/gallery/upload",
		&buf,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func (h *Handler) AddCatalogCategory(c *gin.Context) {
	slug := c.Param("slug")
	raw, _ := io.ReadAll(c.Request.Body)
	h.proxyJSON(c, http.MethodPost, PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/catalog/categories", raw)
}

func (h *Handler) DeleteCatalogCategory(c *gin.Context) {
	slug := c.Param("slug")
	id := c.Param("id")
	h.proxyJSON(c, http.MethodDelete, PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/catalog/categories/"+id, nil)
}

func (h *Handler) AddCatalogItem(c *gin.Context) {
	slug := c.Param("slug")
	raw, _ := io.ReadAll(c.Request.Body)
	h.proxyJSON(c, http.MethodPost, PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/catalog/items", raw)
}

func (h *Handler) DeleteCatalogItem(c *gin.Context) {
	slug := c.Param("slug")
	id := c.Param("id")
	h.proxyJSON(c, http.MethodDelete, PublicAPIBaseURL+"/api/v1/admin/sections/"+slug+"/catalog/items/"+id, nil)
}

// общий помощник
func (h *Handler) proxyJSON(c *gin.Context, method, url string, body []byte) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		c.JSON(500, gin.H{"error": "cant create request"})
		return
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "cant contact api", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}
