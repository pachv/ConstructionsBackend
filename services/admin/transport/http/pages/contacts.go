package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

func (p *Pages) ContactsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/contacts.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	type ContactsPageData struct {
		Base
	}

	data := ContactsPageData{
		Base: p.CreateBase(username, "Контакты", " contacts"),
	}
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
