package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/prayer"
)

type PrayerUpdate struct {
	PrayerID          string `json:"prayerId,omitempty"` // если нужен ID
	PrayerName        string `json:"prayerName"`
	PrayerDescription string `json:"prayerDescription"`
}

// структура для запроса с фронта
type SetPrayersRequest struct {
	Prayers []PrayerUpdate `json:"prayers"`
}

func (h *Handler) SetPrayers(c *gin.Context) {
	var req SetPrayersRequest

	// парсим JSON из body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	var prayersToSend []prayer.PrayerUpdate
	for _, p := range req.Prayers {
		prayersToSend = append(prayersToSend, prayer.PrayerUpdate{
			PrayerID:          p.PrayerID,
			PrayerName:        p.PrayerName,
			PrayerDescription: p.PrayerDescription,
		})
	}

	// вызываем сервис, который отправляет данные на http://prayer:8080/admin/set-prayers
	if err := prayer.SetPrayers(prayersToSend); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prayers: " + err.Error()})
		return
	}

	// успешный ответ
	c.JSON(http.StatusOK, gin.H{"status": "Prayers successfully updated"})
}
