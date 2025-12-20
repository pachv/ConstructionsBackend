package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const constructionsServiceURL = "http://constructions_service:8080/admin/email"

type setAdminEmailRequest struct {
	Email string `json:"email"`
}

func (h *Handler) SetAdminEmailProxy(c *gin.Context) {
	var req setAdminEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest(
		http.MethodPost,
		constructionsServiceURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot create request",
		})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "constructions_service unavailable",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to save email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
