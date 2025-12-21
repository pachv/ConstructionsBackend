package pages

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

type SectionSummary struct {
	ID         string
	Title      string
	Label      string // ✅ ДОБАВИЛИ
	Slug       string
	Image      string
	HasGallery bool
	HasCatalog bool
}

type SectionsPageData struct {
	Base

	Items       []SectionSummary
	Search      string
	CurrentPage int
	PageAmount  int
	Total       int

	InsideBase string

	Error string
}

func (p *Pages) SectionsListPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/sections_list.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	page := 1
	if v := c.Query("page"); v != "" {
		if pv, err := strconv.Atoi(v); err == nil && pv > 0 {
			page = pv
		}
	}
	search := c.Query("search")
	orderBy := c.Query("orderBy") // опционально
	if orderBy == "" {
		orderBy = "title"
	}

	// дергаем admin API через sender
	resp, err := sender.GetAdminSections(c.Request.Context(), page, search, orderBy)
	if err != nil {
		// покажем страницу, но с ошибкой (чтобы шаблон не падал)
		data := SectionsPageData{
			Base:        p.CreateBase(username, "Секции", "sections"),
			Items:       []SectionSummary{},
			InsideBase:  "/admin-service/admin",
			Search:      search,
			CurrentPage: page,
			PageAmount:  1,
			Total:       0,
			Error:       err.Error(),
		}
		if execErr := tmpl.Execute(c.Writer, data); execErr != nil {
			c.String(http.StatusInternalServerError, execErr.Error())
		}
		return
	}

	items := make([]SectionSummary, 0, len(resp.Items))
	for _, it := range resp.Items {
		items = append(items, SectionSummary{
			ID:         it.ID,
			Title:      it.Title,
			Slug:       it.Slug,
			Image:      it.Image,
			HasGallery: it.HasGallery,
			HasCatalog: it.HasCatalog,
		})
	}

	data := SectionsPageData{
		Base:        p.CreateBase(username, "Секции", "sections"),
		Items:       items,
		Search:      search,
		CurrentPage: resp.Page,
		PageAmount:  resp.PageAmount,
		Total:       resp.Total,
		Error:       "",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
