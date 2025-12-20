package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
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

	email, err := sender.GetAdminEmail(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := EmailPageData{
		Base:  p.CreateBase(username, "Email", "email"),
		Email: email,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
