package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/ton"
)

type TonPageData struct {
	Base
	Wallet string
	URL    string

	TonSaveURl string
}

func (p *Pages) TonPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/ton.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tonData, err := ton.GetTonData()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	data := TonPageData{
		Base:       p.CreateBase(username, "Ton", "ton"),
		Wallet:     tonData.Wallet,
		URL:        tonData.URL,
		TonSaveURl: p.Domain + "/admin-service/admin/set-ton",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
