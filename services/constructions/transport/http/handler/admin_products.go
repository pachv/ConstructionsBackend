package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/pachv/constructions/constructions/internal/services"
)

// =========================
// ADMIN CATALOG: CATEGORIES
// =========================

// GET /admin/catalog/categories
// func (h *Handler) AdminGetCatalogCategories(c *gin.Context) {
// 	items, err := h.productAdminService.GetAllCategories()
// 	if err != nil {
// 		h.logger.Error("AdminGetCatalogCategories failed", "err", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, items)
// }

// GET /admin/catalog/categories/by-title/:title
func (h *Handler) AdminGetCatalogCategoryByTitle(c *gin.Context) {
	title := strings.TrimSpace(c.Param("title"))
	item, err := h.productAdminService.GetCategoryByTitle(title)
	if err != nil {
		h.logger.Error("AdminGetCatalogCategoryByTitle failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func boolFromForm(v string) bool {
	v = strings.TrimSpace(strings.ToLower(v))
	return v == "1" || v == "true" || v == "on" || v == "yes"
}

func slugify(s string) string {
	slug.MaxLength = 0 // не резать
	slug.Lowercase = true
	return slug.Make(s)
}

// ✅ POST /admin/catalog/categories - теперь с multipart
func (h *Handler) AdminCreateCatalogCategory(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}

	// ✅ slug ВСЕГДА из title
	slugValue := slugify(title)
	if slugValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot build slug from title"})
		return
	}

	// Файл (опционально)
	var imagePath string
	fh, err := c.FormFile("image")
	if err != nil {
		fh, _ = c.FormFile("imagePath") // для совместимости
	}
	if fh != nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		imagePath, err = h.productAdminService.SaveProductImage(f, filepath.Base(fh.Filename))
		if err != nil {
			h.logger.Error("SaveProductImage failed", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := h.productAdminService.CreateCategory(title, slugValue, imagePath, ""); err != nil {
		h.logger.Error("AdminCreateCatalogCategory failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ✅ PUT /admin/catalog/categories/:id - теперь с multipart
func (h *Handler) AdminUpdateCatalogCategory(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}

	// ✅ slug ВСЕГДА из title
	slugValue := slugify(title)
	if slugValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot build slug from title"})
		return
	}

	// Файл (опционально)
	var imagePath string
	fh, err := c.FormFile("image")
	if err != nil {
		fh, _ = c.FormFile("imagePath")
	}
	if fh != nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		imagePath, err = h.productAdminService.SaveProductImage(f, filepath.Base(fh.Filename))
		if err != nil {
			h.logger.Error("SaveProductImage failed", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := h.productAdminService.UpdateCategory(id, title, slugValue, imagePath); err != nil {
		h.logger.Error("AdminUpdateCatalogCategory failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /admin/catalog/categories/:id
func (h *Handler) AdminDeleteCatalogCategory(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	if err := h.productAdminService.DeleteCategory(id); err != nil {
		h.logger.Error("AdminDeleteCatalogCategory failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// =========================
// ADMIN CATALOG: SECTIONS
// =========================

// GET /admin/catalog/sections
// func (h *Handler) AdminGetCatalogSections(c *gin.Context) {
// 	items, err := h.productAdminService.GetAllSections()
// 	if err != nil {
// 		h.logger.Error("AdminGetCatalogSections failed", "err", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sections"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, items)
// }

// GET /admin/catalog/sections/by-title/:title
func (h *Handler) AdminGetCatalogSectionByTitle(c *gin.Context) {
	title := strings.TrimSpace(c.Param("title"))
	item, err := h.productAdminService.GetSectionByTitle(title)
	if err != nil {
		h.logger.Error("AdminGetCatalogSectionByTitle failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "section not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get section"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// GET /admin/catalog/sections/:id/category
func (h *Handler) AdminGetCatalogSectionCategory(c *gin.Context) {
	sectionID := strings.TrimSpace(c.Param("id"))
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	item, err := h.productAdminService.GetCategoryOfSection(sectionID)
	if err != nil {
		h.logger.Error("AdminGetCatalogSectionCategory failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found for section"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get section category"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// ✅ POST /admin/catalog/sections - теперь с multipart
func (h *Handler) AdminCreateCatalogSection(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	parentCategorySlug := strings.TrimSpace(c.PostForm("parentCategorySlug"))

	if title == "" || parentCategorySlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and parentCategorySlug required"})
		return
	}

	// ✅ slug ВСЕГДА из title
	slugValue := slugify(title)
	if slugValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot build slug from title"})
		return
	}

	// Файл (опционально)
	var imagePath string
	fh, err := c.FormFile("image")
	if err != nil {
		fh, _ = c.FormFile("imagePath")
	}
	if fh != nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		imagePath, err = h.productAdminService.SaveProductImage(f, filepath.Base(fh.Filename))
		if err != nil {
			h.logger.Error("SaveProductImage failed", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := h.productAdminService.CreateSection("", title, slugValue, parentCategorySlug, imagePath); err != nil {
		h.logger.Error("AdminCreateCatalogSection failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ✅ PUT /admin/catalog/sections/:id - теперь с multipart
func (h *Handler) AdminUpdateCatalogSection(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	parentCategorySlug := strings.TrimSpace(c.PostForm("parentCategorySlug"))

	if title == "" || parentCategorySlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and parentCategorySlug required"})
		return
	}

	// ✅ slug ВСЕГДА из title
	slugValue := slugify(title)
	if slugValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot build slug from title"})
		return
	}

	// Файл (опционально)
	var imagePath string
	fh, err := c.FormFile("image")
	if err != nil {
		fh, _ = c.FormFile("imagePath")
	}
	if fh != nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		imagePath, err = h.productAdminService.SaveProductImage(f, filepath.Base(fh.Filename))
		if err != nil {
			h.logger.Error("SaveProductImage failed", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := h.productAdminService.UpdateSection(id, title, slugValue, parentCategorySlug, imagePath); err != nil {
		h.logger.Error("AdminUpdateCatalogSection failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "section not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /admin/catalog/sections/:id
func (h *Handler) AdminDeleteCatalogSection(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	if err := h.productAdminService.DeleteSection(id); err != nil {
		h.logger.Error("AdminDeleteCatalogSection failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete section"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// =========================
// ADMIN CATALOG: PRODUCTS
// =========================

// GET /admin/catalog/products
// func (h *Handler) AdminGetCatalogProducts(c *gin.Context) {
// 	items, err := h.productAdminService.GetAllProducts()
// 	if err != nil {
// 		h.logger.Error("AdminGetCatalogProducts failed", "err", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, items)
// }

// GET /admin/catalog/products/:id
func (h *Handler) AdminGetCatalogProduct(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	item, err := h.productAdminService.GetProduct(id)
	if err != nil {
		h.logger.Error("AdminGetCatalogProduct failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) AdminCreateCatalogProduct(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	categorySlug := strings.TrimSpace(c.PostForm("categorySlug"))
	sectionSlug := strings.TrimSpace(c.PostForm("sectionSlug"))
	brand := strings.TrimSpace(c.PostForm("brand"))
	typ := strings.TrimSpace(c.PostForm("type"))

	priceStr := strings.TrimSpace(c.PostForm("price"))

	if title == "" || categorySlug == "" || sectionSlug == "" || typ == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, categorySlug, sectionSlug, type, price required"})
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	inStock := boolFromForm(c.PostForm("inStock"))

	// discount: просто число процентов
	discount := 0
	discountStr := strings.TrimSpace(c.PostForm("discount"))
	if discountStr != "" {
		v, err := strconv.Atoi(discountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount"})
			return
		}
		discount = v
	}

	// badges[] (и запасной вариант badges)
	badges := c.PostFormArray("badges[]")
	if len(badges) == 0 {
		badges = c.PostFormArray("badges")
	}
	var cleanedBadges []string
	for _, b := range badges {
		b = strings.TrimSpace(b)
		if b != "" {
			cleanedBadges = append(cleanedBadges, b)
		}
	}

	// ✅ slug ВСЕГДА из title (игнорируем пришедший slug)
	slugValue := slugify(title)
	if slugValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot build slug from title"})
		return
	}

	// файл: принимаем "image" и совместимость со старым "imagePath"
	var imagePath string
	fh, err := c.FormFile("image")
	if err != nil {
		fh, _ = c.FormFile("imagePath")
	}
	if fh != nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		// сохраняем файл через сервис (в примере он пишет в ./uploads)
		imagePath, err = h.productAdminService.SaveProductImage(f, filepath.Base(fh.Filename))
		if err != nil {
			h.logger.Error("SaveProductImage failed", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	req := services.CreateProductReq{
		// ID можно не передавать — сервис сгенерит
		Title:        title,
		Slug:         slugValue, // сервис всё равно перезапишет по title
		CategorySlug: categorySlug,
		SectionSlug:  sectionSlug,
		Brand:        brand,
		Type:         typ,
		Price:        price,
		InStock:      inStock,
		ImagePath:    imagePath,
		Badges:       cleanedBadges,
		Discount:     discount,
	}

	if err := h.productAdminService.CreateProduct(req); err != nil {
		h.logger.Error("AdminCreateCatalogProduct failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ✅ PUT /admin/catalog/products/:id - теперь с multipart как создание
func (h *Handler) AdminUpdateCatalogProduct(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	categorySlug := strings.TrimSpace(c.PostForm("categorySlug"))
	sectionSlug := strings.TrimSpace(c.PostForm("sectionSlug"))
	brand := strings.TrimSpace(c.PostForm("brand"))
	typ := strings.TrimSpace(c.PostForm("type"))
	priceStr := strings.TrimSpace(c.PostForm("price"))

	if title == "" || categorySlug == "" || sectionSlug == "" || typ == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, categorySlug, sectionSlug, type, price required"})
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	inStock := boolFromForm(c.PostForm("inStock"))

	// discount: просто число процентов
	discount := 0
	discountStr := strings.TrimSpace(c.PostForm("discount"))
	if discountStr != "" {
		v, err := strconv.Atoi(discountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid discount"})
			return
		}
		discount = v
	}

	// badges[]
	badges := c.PostFormArray("badges[]")
	if len(badges) == 0 {
		badges = c.PostFormArray("badges")
	}
	var cleanedBadges []string
	for _, b := range badges {
		b = strings.TrimSpace(b)
		if b != "" {
			cleanedBadges = append(cleanedBadges, b)
		}
	}

	// ✅ slug ВСЕГДА из title
	slugValue := slugify(title)
	if slugValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot build slug from title"})
		return
	}

	// Файл (опционально)
	var imagePath string
	fh, err := c.FormFile("image")
	if err != nil {
		fh, _ = c.FormFile("imagePath")
	}
	if fh != nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open upload: " + err.Error()})
			return
		}
		defer f.Close()

		imagePath, err = h.productAdminService.SaveProductImage(f, filepath.Base(fh.Filename))
		if err != nil {
			h.logger.Error("SaveProductImage failed", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	req := services.UpdateProductReq{
		Title:        title,
		Slug:         slugValue,
		CategorySlug: categorySlug,
		SectionSlug:  sectionSlug,
		Brand:        brand,
		Type:         typ,
		Price:        price,
		InStock:      inStock,
		ImagePath:    imagePath,
		Badges:       cleanedBadges,
		Discount:     discount,
	}

	if err := h.productAdminService.UpdateProduct(id, req); err != nil {
		h.logger.Error("AdminUpdateCatalogProduct failed", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DELETE /admin/catalog/products/:id
func (h *Handler) AdminDeleteCatalogProduct(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	if err := h.productAdminService.DeleteProduct(id); err != nil {
		h.logger.Error("AdminDeleteCatalogProduct failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /admin/catalog/categories
func (h *Handler) AdminGetCatalogCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	search := c.Query("search")
	orderBy := c.Query("orderBy")

	items, total, err := h.productAdminService.GetAllCategoriesPaginated(page, perPage, search, orderBy)
	if err != nil {
		h.logger.Error("AdminGetCatalogCategories failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":   items,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	})
}

// GET /admin/catalog/sections
func (h *Handler) AdminGetCatalogSections(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	search := c.Query("search")
	orderBy := c.Query("orderBy")

	items, total, err := h.productAdminService.GetAllSectionsPaginated(page, perPage, search, orderBy)
	if err != nil {
		h.logger.Error("AdminGetCatalogSections failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":   items,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	})
}

// GET /admin/catalog/products
func (h *Handler) AdminGetCatalogProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	search := c.Query("search")
	orderBy := c.Query("orderBy")
	categorySlug := c.Query("categorySlug")
	sectionSlug := c.Query("sectionSlug")

	items, total, err := h.productAdminService.GetAllProductsPaginated(page, perPage, search, orderBy, categorySlug, sectionSlug)
	if err != nil {
		h.logger.Error("AdminGetCatalogProducts failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":   items,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	})
}
