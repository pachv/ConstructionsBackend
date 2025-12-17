package handler

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateReview(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	position := strings.TrimSpace(c.PostForm("position"))
	text := strings.TrimSpace(c.PostForm("text"))

	ratingStr := strings.TrimSpace(c.PostForm("rating"))
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating"})
		return
	}

	consentStr := strings.TrimSpace(c.PostForm("consent"))
	consent := consentStr == "true" || consentStr == "1" || strings.EqualFold(consentStr, "yes")

	imagePath := ""

	// фото опционально
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
			h.logger.Error("mkdir uploads error", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		filename := uuid.NewString() + ext
		dst := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			h.logger.Error("save uploaded file error", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		// Если раздаёшь статику так: r.Static("/uploads", "./uploads")
		imagePath = "/uploads/reviews/" + filename
	}

	rv, err := h.reviewService.Create(name, position, text, rating, imagePath, consent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": rv.Id})
}

func (h *Handler) GetPublishedReviews(c *gin.Context) {
	items, err := h.reviewService.GetAllPublished()
	if err != nil {
		h.logger.Error("get reviews error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) GetReviewPicture(c *gin.Context) {
	filename := c.Param("name")

	// защита от ../
	filename = filepath.Base(filename)

	filePath := filepath.Join("./uploads/reviews", filename)

	h.logger.Info("get review picture", "path", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+filename)

	c.File(filePath)
}
