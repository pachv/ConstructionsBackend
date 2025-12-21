package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// const constructionsBaseURL = "http://constructions_service:8080"

// GET /admin-service/admin/sections?...
func (h *Handler) AdminProxyGetSections(c *gin.Context) {
	req, err := http.NewRequest(http.MethodGet, constructionsBaseURL+"/admin/sections?"+c.Request.URL.RawQuery, nil)
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

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// GET /admin-service/admin/sections/:slug
func (h *Handler) AdminProxyGetSectionBySlug(c *gin.Context) {
	slug := c.Param("slug")

	req, err := http.NewRequest(http.MethodGet, constructionsBaseURL+"/admin/sections/"+slug, nil)
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

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// POST /admin-service/admin/sections
func (h *Handler) AdminProxyCreateSection(c *gin.Context) {
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	req, err := http.NewRequest(http.MethodPost, constructionsBaseURL+"/admin/sections", bytes.NewReader(bodyBytes))
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

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

// PUT /admin-service/admin/sections/:id
func (h *Handler) AdminProxyUpdateSection(c *gin.Context) {
	id := c.Param("id")
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	req, err := http.NewRequest(http.MethodPut, constructionsBaseURL+"/admin/sections/"+id, bytes.NewReader(bodyBytes))
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

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

// DELETE /admin-service/admin/sections/:id
func (h *Handler) AdminProxyDeleteSection(c *gin.Context) {
	id := c.Param("id")

	req, err := http.NewRequest(http.MethodDelete, constructionsBaseURL+"/admin/sections/"+id, nil)
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

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}
