package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type EmailPageData struct {
	Base
	Email string
}

func (p *Pages) EmailPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/email.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	// ✅ тестовый email (без сервисов)
	data := EmailPageData{
		Base:  p.CreateBase(username, "Email", "email"),
		Email: "admin@example.com",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
