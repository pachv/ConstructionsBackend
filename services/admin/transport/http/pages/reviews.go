package pages

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

type ReviewMock struct {
	ID         string
	Name       string
	Position   string
	Text       string
	Rating     int
	CanPublish bool
	CreatedAt  time.Time

	PhotoURL string // ✅ добавили
}

type ReviewsMockPageData struct {
	Base

	Items       []ReviewMock
	Search      string
	CurrentPage int
	PageAmount  int
}

func (p *Pages) ReviewsMockPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/reviews.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	page := 1
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	search := c.Query("search")

	resp, err := sender.GetAdminReviews(c.Request.Context(), page, search, "created_at")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]ReviewMock, 0, len(resp.Items))
	for _, it := range resp.Items {
		items = append(items, ReviewMock{
			ID:         it.ID,
			Name:       it.Name,
			Position:   it.Position,
			Text:       it.Text,
			Rating:     it.Rating,
			CanPublish: it.CanPublish,
			CreatedAt:  it.CreatedAt,

			PhotoURL: it.ImagePath,
		})
	}

	data := ReviewsMockPageData{
		Base:        p.CreateBase(username, "Отзывы", "reviews"),
		Items:       items,
		Search:      search,
		CurrentPage: resp.Page,
		PageAmount:  resp.PageAmount,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
