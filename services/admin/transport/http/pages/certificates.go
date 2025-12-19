package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/certificates"
)

type CertificatesPageData struct {
	Base
	Certificates []certificates.Item
}

func (p *Pages) CertificatesPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("./templates/base.html", "./templates/certificates.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	items, err := certificates.GetAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusBadGateway, "failed to fetch certificates: "+err.Error())
		return
	}

	username := c.GetString("username")

	data := CertificatesPageData{
		Base:         p.CreateBase(username, "Сертификаты", "certificates"),
		Certificates: items,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
