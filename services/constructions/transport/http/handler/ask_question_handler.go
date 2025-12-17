package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/services"
)

type AskQuestionRequest struct {
	Message string `json:"message" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Product string `json:"product" binding:"required"`
}

func (h *Handler) AskQuestion(c *gin.Context) {
	var req AskQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	q, err := h.askQuestionService.Create(req.Message, req.Name, req.Email, req.Product)
	if err != nil {
		if err == services.ErrInvalidAskQuestion {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		h.logger.Error("ask question create error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": q.Id})
}
