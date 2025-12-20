package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminDashboard(c *gin.Context) {
	stats, err := h.adminDashboardService.GetStats(c.Request.Context())
	if err != nil {
		h.logger.Error("AdminDashboard: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
