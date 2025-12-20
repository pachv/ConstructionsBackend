package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const constructionsAdminBaseURL = "http://constructions_service:8080/admin"

func (h *Handler) UploadCertificate(c *gin.Context) {
	// читаем из формы
	title := c.PostForm("title")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// собираем multipart заново
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// title (опционально)
	if title != "" {
		_ = writer.WriteField("title", title)
	}

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		_ = writer.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		_ = writer.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to copy file"})
		return
	}

	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, constructionsAdminBaseURL+"/certificates", &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "constructions_service unavailable"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func (h *Handler) DeleteCertificate(c *gin.Context) {
	id := c.Param("id")

	req, err := http.NewRequest("DELETE", "http://constructions_service:8080/admin/certificates/"+id, nil)
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
