package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type CreateUserData struct {
	Base

	CreateUserURL string
	CancelURL     string
}

func (p *Pages) CreateUserPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/create_user.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	data := CreateUserData{
		Base:          p.CreateBase(username, "Create User", "settings"),
		CancelURL:     p.Domain + "/admin/settings/users?page=1",
		CreateUserURL: p.Domain + "/admin-service/admin/create-admin",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
