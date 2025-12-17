package pages

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/referal"
)

type ReferalPages struct {
	Base

	TotalAmount            int
	PremiumPurchasedAmount int

	TotalGift     int
	PurchasedGift int

	SetReferalURL string

	URL         string
	EnglishText string
	ArabText    string
}

func (p *Pages) ReferalPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/referal.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	referalData, err := referal.GetReferalData()
	if err != nil {
		fmt.Println("cant get referal data : " + err.Error())
		panic(err)
	}

	username := c.GetString("username")

	data := ReferalPages{
		Base:                   p.CreateBase(username, "Friends", " friends"),
		TotalAmount:            int(referalData.TotalAmount),
		PremiumPurchasedAmount: int(referalData.PurchasedAmount),
		TotalGift:              int(referalData.TotalGift),
		PurchasedGift:          int(referalData.PurchasedGift),

		SetReferalURL: p.Domain + "/admin-service/admin/set-referal",
		URL:           referalData.URL,
		EnglishText:   referalData.EnglishText,
		ArabText:      referalData.ArabText,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
