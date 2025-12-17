package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Domain string
}

func (p *Pages) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", LoginData{
		Domain: p.Domain,
	})
}
