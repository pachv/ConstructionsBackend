package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// proxyJSON — универсальный прокси для JSON/text ответов
func (h *Handler) proxyJSON(
	c *gin.Context,
	method string,
	url string,
	body []byte,
) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(
		c.Request.Context(),
		method,
		url,
		reader,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// если есть тело — ставим content-type
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact api"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// возвращаем клиенту ровно то, что вернул API
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	c.Data(resp.StatusCode, contentType, respBody)
}
