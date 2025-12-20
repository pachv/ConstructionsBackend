package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pachv/constructions/constructions/internal/services"
)

func (h *Handler) AdminGetReviews(c *gin.Context) {
	page := 1
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	search := strings.TrimSpace(c.Query("search"))
	orderBy := strings.TrimSpace(c.Query("orderBy"))

	items, pageAmount, err := h.adminReviewService.GetPaged(c.Request.Context(), page, search, orderBy)
	if err != nil {
		h.logger.Error("AdminGetReviews: failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      items,
		"page":       page,
		"pageAmount": pageAmount,
	})
}

func (h *Handler) AdminDeleteReview(c *gin.Context) {
	id := c.Param("id")

	// достанем отзыв, чтобы удалить картинку (best-effort)
	rv, err := h.adminReviewService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("AdminDeleteReview: get failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}
	if rv == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if err := h.adminReviewService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("AdminDeleteReview: delete failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	// удалить файл если был
	if strings.TrimSpace(rv.ImagePath) != "" {
		// в базе лежит только имя файла
		name := filepath.Base(rv.ImagePath)
		if name != "" && name != "." {
			_ = os.Remove(filepath.Join("uploads/reviews", name))
		}
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminBulkUpdateReviews(c *gin.Context) {
	var items []services.BulkUpdateReview
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.adminReviewService.BulkUpdate(c.Request.Context(), items); err != nil {
		h.logger.Error("AdminBulkUpdateReviews: failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminCreateReview(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	position := strings.TrimSpace(c.PostForm("position"))
	text := strings.TrimSpace(c.PostForm("text"))

	if name == "" || position == "" || text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, position, text are required"})
		return
	}

	ratingStr := strings.TrimSpace(c.PostForm("rating"))
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating"})
		return
	}

	consentStr := strings.TrimSpace(c.PostForm("consent"))
	consent := consentStr == "true" || consentStr == "1" || strings.EqualFold(consentStr, "yes")

	imagePath := ""

	// photo optional
	file, err := c.FormFile("photo")
	if err == nil && file != nil {
		ct := file.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "photo must be an image"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == "" {
			ext = ".jpg"
		}

		uploadDir := "uploads/reviews"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			h.logger.Error("AdminCreateReview: mkdir failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		filename := uuid.NewString() + ext
		dst := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			h.logger.Error("AdminCreateReview: save file failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		imagePath = filename // в базу кладём только имя
	}

	id, err := h.adminReviewService.Create(c.Request.Context(), name, position, text, rating, imagePath, consent)
	if err != nil {
		h.logger.Error("AdminCreateReview: failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
