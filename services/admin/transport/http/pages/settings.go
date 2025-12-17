package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type SettingsData struct {
	Base

	UsersPageURL string
}

func (p *Pages) SettingsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/settings.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	data := SettingsData{
		Base:         p.CreateBase(username, "Settings", " settings"),
		UsersPageURL: p.Domain + "/admin/settings/users?page=1",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
