package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/prices"
)

func (h *Handler) UpdatePrices(c *gin.Context) {
	var pricesToUpdate []prices.UpdatePrice

	// Пробуем распарсить массив из JSON
	if err := c.ShouldBindJSON(&pricesToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Вызываем твою функцию отправки
	if err := prices.UpdatePrices(pricesToUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update prices",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "prices updated successfully"})
}
