package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/responses"
)

func (h *Handler) InitInsideHandlers(r *gin.RouterGroup) {
	r.GET("/favicon", h.GetFavicon)
	r.GET("/logo", h.GetLogo)
	r.POST("/logout", h.LogoutUser)

	r.POST("/set-referal", h.SetReferal)
	r.POST("/set-ton", h.SetTon)
	r.POST("/set-translationss", h.SetTransaltionsService)
	r.POST("/set-prayers", h.SetPrayers)

	r.POST("/update-prices", h.UpdatePrices)

	r.POST("/create-admin", h.CreateAdmin)
	r.POST("/delete-admin/:id", h.DeleteAdmin)

	r.POST("/update-admin", h.UpdateAdmin)
	r.POST("/update-bot-data", h.UpdateTextHandler)
	r.POST("/upload-bot-img", h.UploadBotImage)

	// email

	r.POST("/set-admin-email", h.SetAdminEmailProxy)

	// certificates
	r.POST("/certificates", h.UploadCertificate)
	r.DELETE("/certificates/:id", h.DeleteCertificate)

	// sections

	r.GET("/sections", h.ProxySectionsList)        // -> GET {API}/api/v1/sections
	r.GET("/sections/:slug", h.ProxySectionBySlug) // -> GET {API}/api/v1/sections/:slug

	// ✅ gallery
	r.POST("/sections/:slug/gallery", h.AddGalleryItem)              // -> POST {API}/api/v1/admin/sections/:slug/gallery
	r.DELETE("/sections/:slug/gallery/:id", h.DeleteGalleryItem)     // -> DELETE {API}/api/v1/admin/sections/:slug/gallery/:id
	r.POST("/sections/:slug/gallery/upload", h.UploadGalleryPicture) // -> POST {API}/api/v1/admin/sections/:slug/gallery/upload (multipart)

	// ✅ catalog categories
	r.POST("/sections/:slug/catalog/categories", h.AddCatalogCategory)
	r.DELETE("/sections/:slug/catalog/categories/:id", h.DeleteCatalogCategory)

	// ✅ catalog items
	r.POST("/sections/:slug/catalog/items", h.AddCatalogItem)
	r.DELETE("/sections/:slug/catalog/items/:id", h.DeleteCatalogItem)
}

// Handler: принимает файл и пересылает его на микросервис бота
func (h *Handler) UploadBotImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Готовим multipart-запрос
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form file"})
		return
	}
	if _, err := io.Copy(part, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to copy file"})
		return
	}
	writer.Close()

	// Отправляем файл на микросервис бота
	req, err := http.NewRequest("POST", "http://bot:8080/set-image", &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send request to bot"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image uploaded successfully"})
}

type UpdateTextData struct {
	WelcomeText          string `db:"welcome_text" json:"welcomeText"`
	UnknownText          string `db:"unknown_text" json:"unknownText"`
	ReferalActivatedText string `db:"referal_activated_text" json:"referalAtivated"`
	AppUrl               string `db:"app_url" json:"appURL"`
	PrayerEngText        string `db:"prayer_eng_text" json:"prayerEngText"`
	PrayerArText         string `db:"prayer_ar_text" json:"prayerArText"`
}

func (h *Handler) UpdateTextHandler(c *gin.Context) {
	var data UpdateTextData

	// 🧩 Читаем JSON из фронта
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"details": err.Error(),
		})
		return
	}

	// 🔄 Конвертируем обратно в JSON для пересылки
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to serialize request data",
		})
		return
	}

	// 🌐 Отправляем запрос на bot-сервис
	resp, err := http.Post("http://bot:8080/set-bot-data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "Failed to contact bot service",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 📦 Читаем ответ от bot
	body, _ := io.ReadAll(resp.Body)

	// 🔁 Проксируем статус и тело обратно фронту
	c.Data(resp.StatusCode, "application/json", body)
}

type CreateAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) CreateAdmin(c *gin.Context) {

	var req CreateAdminRequest

	if err := c.BindJSON(&req); err != nil {
		responses.BadRequest(c, "cant bind json")
		c.Abort()
		return
	}

	err := h.authService.CreateUser(req.Username, req.Password)
	if err != nil {
		responses.BadRequest(c, "cant bind json")
		c.Abort()
		return
	}

	responses.Ok(c, gin.H{})
}

func (h *Handler) DeleteAdmin(c *gin.Context) {

	id := c.Param("id")

	err := h.authService.DeleteUser(id)
	if err != nil {
		responses.BadRequest(c, "cant delete user")
		c.Abort()
		return
	}

	responses.Ok(c, "")
}

type UpdateUserRequest struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) UpdateAdmin(c *gin.Context) {

	var req UpdateUserRequest

	if err := c.BindJSON(&req); err != nil {
		responses.BadRequest(c, "cant bind json : "+err.Error())
		c.Abort()
		return
	}

	err := h.authService.UpdateUser(req.Id, req.Username, req.Password)
	if err != nil {
		responses.BadRequest(c, "cant bind json : "+err.Error())
		c.Abort()
		return
	}

	responses.Ok(c, gin.H{})
}
