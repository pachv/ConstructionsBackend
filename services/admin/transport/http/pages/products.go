package pages

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

type ProductCard struct {
	ID        string
	Title     string
	Slug      string
	Category  string
	Section   string
	Brand     string
	Type      string
	PriceRub  int
	InStock   bool
	SaleValue int
	BadgesCSV string
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

	CategoriesJSON string
	SectionsJSON   string

	Pager Pager

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

	// --------------------
	// Query
	// --------------------
	page := 1
	if v := strings.TrimSpace(c.Query("page")); v != "" {
		if pv, err := strconv.Atoi(v); err == nil && pv > 0 {
			page = pv
		}
	}

	activeTab := strings.TrimSpace(c.Query("tab"))
	if activeTab == "" {
		activeTab = "products"
	}

	search := strings.TrimSpace(c.Query("search"))
	orderBy := strings.TrimSpace(c.Query("orderBy"))
	if orderBy == "" {
		orderBy = "created_at"
	}

	// optional filters (можно использовать позже)
	categorySlug := strings.TrimSpace(c.Query("categorySlug"))
	sectionSlug := strings.TrimSpace(c.Query("sectionSlug"))

	// --------------------
	// Preload categories + sections для селектов в модалках
	// (даже если активная вкладка = products)
	// --------------------
	preloadCatsResp, err := sender.ConstructionsGetCatalogCategories(c.Request.Context(), 1, "", "created_at")
	if err != nil {
		c.String(http.StatusInternalServerError, "Categories preload error: "+err.Error())
		return
	}
	preloadSecsResp, err := sender.ConstructionsGetCatalogSections(c.Request.Context(), 1, "", "created_at")
	if err != nil {
		c.String(http.StatusInternalServerError, "Sections preload error: "+err.Error())
		return
	}

	preloadCategories := make([]CategoryCard, 0, len(preloadCatsResp.Items))
	for _, cat := range preloadCatsResp.Items {
		preloadCategories = append(preloadCategories, CategoryCard{
			ID:       cat.ID,
			Title:    cat.Title,
			Slug:     cat.Slug,
			ImageURL: cat.ImagePath,
		})
	}

	preloadSections := make([]SectionCard, 0, len(preloadSecsResp.Items))
	for _, sec := range preloadSecsResp.Items {
		preloadSections = append(preloadSections, SectionCard{
			ID:                 sec.ID,
			Title:              sec.Title,
			Slug:               sec.Slug,
			ParentCategorySlug: sec.ParentCategorySlug,
			ImageURL:           sec.ImagePath,
		})
	}

	data := ProductsPageData{
		Base:      p.CreateBase(username, "Каталог", "products"),
		ActiveTab: activeTab,
		Search:    search,
		OrderBy:   orderBy,

		Products:   []ProductCard{},
		Categories: preloadCategories,
		Sections:   preloadSections,

		// JSON для JS (селекты)
		CategoriesJSON: categoriesToJSON(preloadCategories),
		SectionsJSON:   sectionsToJSON(preloadSections),

		// urls (как в твоём коде)
		SecCreateURL: "/admin-service/admin/products/sections/create",
		SecUpdateURL: "/admin-service/admin/products/sections/update",
		SecDeleteURL: "/admin-service/admin/products/sections/delete",

		CatCreateURL: "/admin-service/admin/products/categories/create",
		CatUpdateURL: "/admin-service/admin/products/categories/update",
		CatDeleteURL: "/admin-service/admin/products/categories/delete",

		ProdCreateURL: "/admin-service/admin/products/products/create",
		ProdUpdateURL: "/admin-service/admin/products/products/update",
		ProdDeleteURL: "/admin-service/admin/products/products/delete",

		// обязательно, чтобы шаблон не падал
		Pager: buildPager(c, page, 1),
	}

	// --------------------
	// Load active tab + pager
	// --------------------
	switch activeTab {
	case "categories":
		resp, err := sender.ConstructionsGetCatalogCategories(c.Request.Context(), page, search, orderBy)
		if err != nil {
			c.String(http.StatusInternalServerError, "Categories error: "+err.Error())
			return
		}

		// отображаем именно текущую страницу
		items := make([]CategoryCard, 0, len(resp.Items))
		for _, cat := range resp.Items {
			items = append(items, CategoryCard{
				ID:       cat.ID,
				Title:    cat.Title,
				Slug:     cat.Slug,
				ImageURL: cat.ImagePath,
			})
		}
		data.Categories = items

		// Важно: селекты в модалках лучше питать ПРЕЛОАДОМ, а не page=1 текущего поиска.
		// Поэтому CategoriesJSON оставляем preload-версией, чтобы в селектах было "что-то" всегда.
		// Если хочешь, чтобы JSON совпадал с видимым списком — раскомментируй строку ниже:
		// data.CategoriesJSON = categoriesToJSON(items)

		pageAmount := calcPageAmount(resp.Total, resp.PerPage)
		data.Pager = buildPager(c, resp.Page, pageAmount)

	case "sections":
		resp, err := sender.ConstructionsGetCatalogSections(c.Request.Context(), page, search, orderBy)
		if err != nil {
			c.String(http.StatusInternalServerError, "Sections error: "+err.Error())
			return
		}

		items := make([]SectionCard, 0, len(resp.Items))
		for _, sec := range resp.Items {
			items = append(items, SectionCard{
				ID:                 sec.ID,
				Title:              sec.Title,
				Slug:               sec.Slug,
				ParentCategorySlug: sec.ParentCategorySlug,
				ImageURL:           sec.ImagePath,
			})
		}
		data.Sections = items

		// Аналогично: JSON для селектов оставляем preload-версией
		// data.SectionsJSON = sectionsToJSON(items)

		pageAmount := calcPageAmount(resp.Total, resp.PerPage)
		data.Pager = buildPager(c, resp.Page, pageAmount)

	default: // "products"
		resp, err := sender.ConstructionsGetCatalogProducts(c.Request.Context(), page, search, orderBy, categorySlug, sectionSlug)
		if err != nil {
			c.String(http.StatusInternalServerError, "Products error: "+err.Error())
			return
		}

		items := make([]ProductCard, 0, len(resp.Items))
		for _, pr := range resp.Items {
			items = append(items, ProductCard{
				ID:        pr.ID,
				Title:     pr.Title,
				Slug:      pr.Slug,
				Category:  pr.CategorySlug,
				Section:   pr.SectionSlug,
				Brand:     pr.Brand,
				Type:      pr.Type,
				PriceRub:  pr.Price,
				InStock:   pr.InStock,
				SaleValue: pr.SalePercent,
				BadgesCSV: strings.Join(pr.Badges, ","),
				ImageURL:  pr.ImagePath,
			})
		}
		data.Products = items

		pageAmount := calcPageAmount(resp.Total, resp.PerPage)
		data.Pager = buildPager(c, resp.Page, pageAmount)
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func calcPageAmount(total, perPage int) int {
	if perPage <= 0 {
		perPage = 20
	}
	n := (total + perPage - 1) / perPage
	if n < 1 {
		n = 1
	}
	return n
}

// buildPager делает ссылки, сохраняя query-параметры (кроме page, он будет заменён).
func buildPager(c *gin.Context, page, pageAmount int) Pager {
	if page < 1 {
		page = 1
	}
	if pageAmount < 1 {
		pageAmount = 1
	}
	if page > pageAmount {
		page = pageAmount
	}

	u := &url.URL{Path: c.Request.URL.Path}
	q := c.Request.URL.Query()

	p := Pager{
		Page:       page,
		PageAmount: pageAmount,
		HasPrev:    page > 1,
		HasNext:    page < pageAmount,
		PrevURL:    "",
		NextURL:    "",
	}

	if p.HasPrev {
		q.Set("page", strconv.Itoa(page-1))
		u.RawQuery = q.Encode()
		p.PrevURL = u.String()
	}

	if p.HasNext {
		q.Set("page", strconv.Itoa(page+1))
		u.RawQuery = q.Encode()
		p.NextURL = u.String()
	}

	return p
}

// Эти две функции нужны, чтобы безопасно отдать JSON внутрь шаблона,
// а потом распарсить в JS через JSON.parse(`...`)

func categoriesToJSON(items []CategoryCard) string {
	type row struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
	}
	out := make([]row, 0, len(items))
	for _, it := range items {
		out = append(out, row{
			ID:    it.ID,
			Title: it.Title,
			Slug:  it.Slug,
		})
	}
	b, _ := json.Marshal(out)

	// важно: чтобы строка не ломала шаблон/JS (внутри у тебя JSON.parse(`{{.CategoriesJSON}}`))
	s := string(b)
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "</", "<\\/") // защита от </script>
	return s
}

func sectionsToJSON(items []SectionCard) string {
	type row struct {
		ID             string `json:"id"`
		Title          string `json:"title"`
		Slug           string `json:"slug"`
		ParentCategory string `json:"parentCategory"`
	}
	out := make([]row, 0, len(items))
	for _, it := range items {
		out = append(out, row{
			ID:             it.ID,
			Title:          it.Title,
			Slug:           it.Slug,
			ParentCategory: it.ParentCategorySlug, // JS ждёт sec.parentCategory
		})
	}
	b, _ := json.Marshal(out)

	s := string(b)
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "</", "<\\/")
	return s
}
