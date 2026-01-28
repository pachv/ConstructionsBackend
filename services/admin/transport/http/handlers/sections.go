package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const PublicAPIBaseURL = "http://localhost:80/admin-serivce/"

// список (НОВЫЙ)
func (h *Handler) ProxySectionsList(c *gin.Context) {
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURL+"/admin/sections/all", nil)
}

// секция FULL (НОВЫЙ)
func (h *Handler) ProxySectionBySlug(c *gin.Context) {
	slug := c.Param("slug")
	h.proxyJSON(c, http.MethodGet, PublicAPIBaseURL+"/admin/sections/"+slug+"/full", nil)
}

// CREATE basic (multipart): POST /inside/sections -> POST {API}/admin/sections/create-form
func (h *Handler) CreateSectionBasicProxy(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image is required"})
		return
	}

	// собираем multipart заново и прокидываем 1:1
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("title", c.PostForm("title"))
	_ = w.WriteField("advantegesText", c.PostForm("advantegesText"))

	// advanteges[] (много)
	advs := c.PostFormArray("advanteges[]")
	if len(advs) == 0 {
		advs = c.PostFormArray("advanteges")
	}
	for _, a := range advs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		_ = w.WriteField("advanteges[]", a)
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	part, err := w.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(500, gin.H{"error": "failed to copy"})
		return
	}

	_ = w.Close()

	req, err := http.NewRequest(http.MethodPost, PublicAPIBaseURL+"/admin/sections/create-form", &buf)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": "failed to contact api"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// DELETE секции: DELETE /inside/sections/:id -> DELETE {API}/admin/sections/:id
func (h *Handler) DeleteSectionProxy(c *gin.Context) {
	id := c.Param("id")
	h.proxyJSON(c, http.MethodDelete, PublicAPIBaseURL+"/admin/sections/"+id, nil)
}

// update gallery

// update products

// add galery

// delete gallery

// add catgeory

// delete category
