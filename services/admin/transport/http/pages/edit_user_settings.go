package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type UserEditSettings struct {
	Base

	Id       string
	Username string

	UpdateUserURL string
	CancelURL     string
}

func (p *Pages) EditUserPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/edit_user.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	userId := c.Param("id")

	username := c.GetString("username")

	user, err := p.authService.GetUser(userId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := UserEditSettings{
		Base:     p.CreateBase(username, "Settings", "settings"),
		Id:       user.Id,
		Username: user.Username,

		CancelURL:     p.Domain + "/admin/settings/users?page=1",
		UpdateUserURL: p.Domain + "/admin-service/admin/update-admin",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
