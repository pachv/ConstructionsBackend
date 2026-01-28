package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

// -------------------------
// view models
// -------------------------

type ProductCard struct {
	ID       string
	Title    string
	Slug     string
	Category string
	Section  string
	Brand    string
	Type     string
	PriceRub int
	InStock  bool

	SaleValue string // sale_20
	BadgesCSV string // hit,recommended
	ImageURL  string
}

type CategoryCard struct {
	ID       string
	Title    string
	Slug     string
	ImageURL string
}

type SectionCard struct {
	ID                 string
	Title              string
	Slug               string
	ParentCategorySlug string
	ImageURL           string
}

// JSON-структуры для JavaScript
type CategoryJSON struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type SectionJSON struct {
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	ParentCategory string `json:"parentCategory"`
}

type Pager struct {
	Page       int
	PageAmount int
	HasPrev    bool
	HasNext    bool
	PrevURL    string
	NextURL    string
}

type ProductsPageData struct {
	Base

	ActiveTab string
	Search    string
	OrderBy   string

	Products   []ProductCard
	Categories []CategoryCard
	Sections   []SectionCard

	// JSON для JavaScript (селекторов)
	CategoriesJSON template.JS
	SectionsJSON   template.JS

	Pager Pager

	// admin-service handlers urls (как в InitInsideHandlers)
	ProdCreateURL string
	ProdUpdateURL string
	ProdDeleteURL string

	CatCreateURL string
	CatUpdateURL string
	CatDeleteURL string

	SecCreateURL string
	SecUpdateURL string
	SecDeleteURL string
}

// -------------------------
// helpers
// -------------------------

func getTab(c *gin.Context) string {
	tab := strings.TrimSpace(c.Query("tab"))
	switch tab {
	case "products", "categories", "sections":
		return tab
	default:
		return "products"
	}
}

func parsePage(c *gin.Context) int {
	page := 1
	if p := strings.TrimSpace(c.Query("page")); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	return page
}

func sanitizeOrderBy(v string) string {
	v = strings.TrimSpace(v)
	switch v {
	case "created_at", "title", "price":
		return v
	default:
		return "created_at"
	}
}

func buildPublicImageURL(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return ""
	}
	return "http://localhost:80/api/v1/products/picture/" + filename
}

func buildPager(c *gin.Context, page, pageAmount int) Pager {
	if page < 1 {
		page = 1
	}
	if pageAmount < 1 {
		pageAmount = 1
	}

	hasPrev := page > 1
	hasNext := page < pageAmount

	qPrev := c.Request.URL.Query()
	qPrev.Set("page", strconv.Itoa(page-1))
	prev := c.Request.URL.Path + "?" + qPrev.Encode()

	qNext := c.Request.URL.Query()
	qNext.Set("page", strconv.Itoa(page+1))
	next := c.Request.URL.Path + "?" + qNext.Encode()

	return Pager{
		Page:       page,
		PageAmount: pageAmount,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
		PrevURL:    prev,
		NextURL:    next,
	}
}

// Преобразование в JSON для шаблона
func categoriesToJSON(categories []CategoryCard) template.JS {
	items := make([]CategoryJSON, 0, len(categories))
	for _, cat := range categories {
		items = append(items, CategoryJSON{
			Slug:  cat.Slug,
			Title: cat.Title,
		})
	}

	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return template.JS("[]")
	}
	return template.JS(jsonBytes)
}

func sectionsToJSON(sections []SectionCard) template.JS {
	items := make([]SectionJSON, 0, len(sections))
	for _, sec := range sections {
		items = append(items, SectionJSON{
			Slug:           sec.Slug,
			Title:          sec.Title,
			ParentCategory: sec.ParentCategorySlug,
		})
	}

	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return template.JS("[]")
	}
	return template.JS(jsonBytes)
}

// -------------------------
// page
// -------------------------

func (p *Pages) ProductsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/products.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	activeTab := getTab(c)
	page := parsePage(c)

	search := strings.TrimSpace(c.Query("search"))
	orderBy := sanitizeOrderBy(c.Query("orderBy"))

	// ✅ GET списки идут в constructions_service:8080 (через sender.Constructions...)
	productsDTO, err := sender.ConstructionsGetCatalogProducts(
		c.Request.Context(),
		page,
		search,
		orderBy,
		"", // categorySlug filter (если нужно — передай из query)
		"", // sectionSlug filter (если нужно — передай из query)
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	categoriesDTO, err := sender.ConstructionsGetCatalogCategories(c.Request.Context(), page, search, orderBy)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	sectionsDTO, err := sender.ConstructionsGetCatalogSections(c.Request.Context(), page, search, orderBy)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// map products
	products := make([]ProductCard, 0, len(productsDTO))
	for _, pr := range productsDTO {
		sale := ""
		if pr.SalePercent > 0 {
			sale = "sale_" + strconv.Itoa(pr.SalePercent)
		}

		img := ""
		if strings.TrimSpace(pr.ImagePath) != "" {
			img = buildPublicImageURL(pr.ImagePath)
		}

		products = append(products, ProductCard{
			ID:        pr.ID,
			Title:     pr.Title,
			Slug:      pr.Slug,
			Category:  pr.CategorySlug,
			Section:   pr.SectionSlug,
			Brand:     pr.Brand,
			Type:      pr.Type,
			PriceRub:  pr.Price,
			InStock:   pr.InStock,
			SaleValue: sale,
			BadgesCSV: strings.Join(pr.Badges, ","),
			ImageURL:  img,
		})
	}

	categories := make([]CategoryCard, 0, len(categoriesDTO))
	for _, cat := range categoriesDTO {
		img := ""
		if strings.TrimSpace(cat.ImagePath) != "" {
			img = buildPublicImageURL(cat.ImagePath)
		}
		categories = append(categories, CategoryCard{
			ID:       cat.ID,
			Title:    cat.Title,
			Slug:     cat.Slug,
			ImageURL: img,
		})
	}

	sections := make([]SectionCard, 0, len(sectionsDTO))
	for _, sec := range sectionsDTO {
		img := ""
		if strings.TrimSpace(sec.ImagePath) != "" {
			img = buildPublicImageURL(sec.ImagePath)
		}
		sections = append(sections, SectionCard{
			ID:                 sec.ID,
			Title:              sec.Title,
			Slug:               sec.Slug,
			ParentCategorySlug: sec.ParentCategorySlug,
			ImageURL:           img,
		})
	}

	// Преобразуем в JSON для использования в JavaScript
	categoriesJSON := categoriesToJSON(categories)
	sectionsJSON := sectionsToJSON(sections)

	// pager: пока sender отдаёт массив — значит pageAmount = 1
	// (если позже constructions начнёт отдавать объект {items,pageAmount,...} — расширим sender)
	activePages := 1
	pager := buildPager(c, page, activePages)

	fmt.Println(products)

	data := ProductsPageData{
		Base:       p.CreateBase(username, "Каталог", "products"),
		ActiveTab:  activeTab,
		Search:     search,
		OrderBy:    orderBy,
		Products:   products,
		Categories: categories,
		Sections:   sections,

		// ✅ JSON-данные для JavaScript (селекторов в модалках)
		CategoriesJSON: categoriesJSON,
		SectionsJSON:   sectionsJSON,

		Pager: pager,

		// ✅ CRUD идёт в admin-service handlers (как у тебя в InitInsideHandlers)
		SecCreateURL: "/admin-service/admin/products/sections/create",
		SecUpdateURL: "/admin-service/admin/products/sections/update",
		SecDeleteURL: "/admin-service/admin/products/sections/delete",

		CatCreateURL: "/admin-service/admin/products/categories/create",
		CatUpdateURL: "/admin-service/admin/products/categories/update",
		CatDeleteURL: "/admin-service/admin/products/categories/delete",

		ProdCreateURL: "/admin-service/admin/products/products/create",
		ProdUpdateURL: "/admin-service/admin/products/products/update",
		ProdDeleteURL: "/admin-service/admin/products/products/delete",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
