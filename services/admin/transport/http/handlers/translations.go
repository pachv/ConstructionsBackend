package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/translations"
)

type FrontendTranslationsRequest struct {
	Translations string `json:"translations"`
}

func (h *Handler) SetTransaltionsService(c *gin.Context) {
	var req FrontendTranslationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("invalid JSON: %v", err),
		})
		return
	}

	err := translations.SetTranslations(req.Translations)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("failed to send to settings service: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
