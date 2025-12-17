package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AskQuestionRequest struct {
	Message string  `json:"message" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	Phone   string  `json:"phone" binding:"required"`
	Email   *string `json:"email"`   // optional
	Product *string `json:"product"` // optional
	Consent bool    `json:"consent" binding:"required"`
}

func (h *Handler) AskQuestion(c *gin.Context) {
	var req AskQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	q, err := h.askQuestionService.Create(
		req.Message,
		req.Name,
		req.Phone,
		req.Email,
		req.Product,
		req.Consent,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": q.Id})
}
