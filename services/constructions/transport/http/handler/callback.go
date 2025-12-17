package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CallbackRequestDTO struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Consent bool   `json:"consent" binding:"required"`
}

func (h *Handler) Callback(c *gin.Context) {
	var req CallbackRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	cb, err := h.callbackService.Create(req.Name, req.Phone, req.Consent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": cb.Id})
}
