package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

type FontsPageData struct {
	Base
	Fonts []sender.AdminFont
}

func (p *Pages) FontsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/fonts.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	items, err := sender.GetAdminFonts(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if items == nil {
		items = []sender.AdminFont{}
	}

	data := FontsPageData{
		Base:  p.CreateBase(username, "Fonts", "fonts"),
		Fonts: items,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
