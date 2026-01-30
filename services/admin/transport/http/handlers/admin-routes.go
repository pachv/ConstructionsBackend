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

	// gellery

	// ✅ NEW gallery categories/photos (constructions gallery)
	r.POST("/gallery/categories", h.GalleryCreateCategoryProxy)
	r.PUT("/gallery/categories/:id", h.GalleryUpdateCategoryProxy)
	r.DELETE("/gallery/categories/:id", h.GalleryDeleteCategoryProxy)

	r.POST("/gallery/categories/:id/photos", h.GalleryAddPhotoProxy) // multipart: photo
	r.DELETE("/gallery/photos/:id", h.GalleryDeletePhotoProxy)

	// sections

	insideSections := r.Group("/sections")
	{
		// LIST
		insideSections.GET("/all", h.ProxySectionsList)

		// CREATE
		insideSections.POST("/create", h.CreateSectionBasicProxy)

		// UPDATE/DELETE by ID
		insideSections.PUT("/update/:id", h.UpdateSectionProxy)
		insideSections.DELETE("/delete/:id", h.DeleteSectionProxy)

		insideSections.PATCH("/:id/gallery/toggle", h.ToggleSectionGalleryProxy)
		insideSections.PATCH("/:id/catalog/toggle", h.ToggleSectionCatalogProxy)

		// GET by slug
		insideSections.GET("/view/:slug", h.ProxySectionBySlug)

		// GALLERY operations
		insideSections.POST("/gallery/:slug/add", h.AddGalleryPhotoProxy)
		insideSections.DELETE("/gallery/:slug/photo/:photoId", h.DeleteGalleryPhotoProxy)
		insideSections.POST("/gallery/:slug/upload", h.UploadGalleryImageProxy)

		// CATALOG - Categories
		insideSections.POST("/catalog/:slug/categories/add", h.AddCatalogCategoryProxy)
		insideSections.DELETE("/catalog/:slug/categories/:categoryId", h.DeleteCatalogCategoryProxy)
		insideSections.PUT("/catalog/:slug/update", h.UpdateCatalogProxy)

		// CATALOG - Items
		insideSections.POST("/catalog/:slug/items/add", h.AddCatalogItemProxy)
		insideSections.DELETE("/catalog/:slug/items/:itemId", h.DeleteCatalogItemProxy)
		insideSections.POST("/catalog/:slug/items/upload", h.UploadCatalogItemImageProxy)
	}
	// Revies
	r.DELETE("/reviews/:id", h.DeleteReviewProxy)
	r.PUT("/reviews/bulk", h.BulkUpdateReviewsProxy)
	r.POST("/reviews", h.CreateReviewProxy)

	// contacts

	r.GET("/contacts", h.ProxyAdminContactsGet)                       // -> GET {API}/admin/contacts
	r.PUT("/contacts/email", h.ProxyAdminContactsSetEmail)            // -> PUT {API}/admin/contacts/email
	r.GET("/contacts/phones", h.ProxyAdminContactsListPhones)         // -> GET {API}/admin/contacts/phones
	r.PUT("/contacts/phones", h.ProxyAdminContactsUpsertPhone)        // -> PUT {API}/admin/contacts/phones
	r.DELETE("/contacts/phones/:id", h.ProxyAdminContactsDeletePhone) // -> DELETE {API}/admin/contacts/phones/:id
	r.PUT("/contacts/phones/reorder", h.ProxyAdminContactsReorderPhones)

	r.GET("/contacts/addresses", h.ProxyAdminContactsListAddresses)            // -> GET {API}/admin/contacts/addresses
	r.PUT("/contacts/addresses", h.ProxyAdminContactsUpsertAddress)            // -> PUT {API}/admin/contacts/addresses
	r.DELETE("/contacts/addresses/:id", h.ProxyAdminContactsDeleteAddress)     // -> DELETE {API}/admin/contacts/addresses/:id
	r.PUT("/contacts/addresses/reorder", h.ProxyAdminContactsReorderAddresses) // -> PUT {API}/admin/contacts/addresses/reorder

	// products
	r.POST("/products/sections/create", h.CreateCatalogSectionFromPage)
	r.POST("/products/sections/update", h.UpdateCatalogSectionFromPage)
	r.POST("/products/sections/delete", h.DeleteCatalogSectionFromPage)

	// добавь аналогично если нужно:
	r.POST("/products/categories/create", h.CreateCatalogCategoryFromPage)
	r.POST("/products/categories/update", h.UpdateCatalogCategoryFromPage)
	r.POST("/products/categories/delete", h.DeleteCatalogCategoryFromPage)

	r.POST("/products/products/create", h.CreateCatalogProductFromPage)
	r.POST("/products/products/update", h.UpdateCatalogProductFromPage)
	r.POST("/products/products/delete", h.DeleteCatalogProductFromPage)

	insideFonts := r.Group("/fonts")
	{
		insideFonts.GET("", h.ProxyAdminFontsList)
		insideFonts.POST("", h.ProxyAdminFontsCreate)
		insideFonts.POST("/:id/select", h.ProxyAdminFontsSelect)
		insideFonts.DELETE("/:id", h.ProxyAdminFontsDelete)
	}

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
