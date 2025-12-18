package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/services"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var req services.CreateOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json", "details": err.Error()})
		return
	}

	// если есть авторизация — достань userID из контекста (пример)
	var userID *string
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok && s != "" {
			userID = &s
		}
	}

	orderID, err := h.orderService.CreateOrder(c.Request.Context(), userID, req, "./templates/order.html")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"orderId": orderID})
}
