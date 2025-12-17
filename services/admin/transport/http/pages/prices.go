package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/prices"
)

type Price struct {
	PriceId           string
	PriceTitle        string
	PriceDescription  string
	PriceRewardType   string
	PriceRewardAmount int
	PriceTonId        string
	PriceStarsId      string
	PriceTonAmount    float64
	PriceStarsAmount  int
}

type PricesPageData struct {
	Base

	EnglishPrices []Price
	ArabicPrices  []Price

	SendPricesURL string
}

func (p *Pages) PricesPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/prices.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	engPrices := []Price{}
	arabPrices := []Price{}

	prices, err := prices.GetPrices()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for _, p := range prices {

		if p.LanguageId == ArabLanguageId {
			arabPrices = append(arabPrices, Price{
				PriceId:           p.PriceId,
				PriceTitle:        p.PriceTitle,
				PriceDescription:  p.PriceDescription,
				PriceRewardType:   p.PriceRewardType,
				PriceRewardAmount: p.PriceRewardAmount,
				PriceTonId:        p.PriceTonId,
				PriceStarsId:      p.PriceStartsId,
				PriceTonAmount:    p.PriceTonAmount,
				PriceStarsAmount:  p.PriceStarsAmount,
			})
		} else {
			engPrices = append(engPrices, Price{
				PriceId:           p.PriceId,
				PriceTitle:        p.PriceTitle,
				PriceDescription:  p.PriceDescription,
				PriceRewardType:   p.PriceRewardType,
				PriceRewardAmount: p.PriceRewardAmount,
				PriceTonId:        p.PriceTonId,
				PriceStarsId:      p.PriceStartsId,
				PriceTonAmount:    p.PriceTonAmount,
				PriceStarsAmount:  p.PriceStarsAmount,
			})
		}
	}

	username := c.GetString("username")

	data := PricesPageData{
		Base:          p.CreateBase(username, "Prices", " prices"),
		EnglishPrices: engPrices,
		ArabicPrices:  arabPrices,
		SendPricesURL: p.Domain + "/admin-service/admin/update-prices",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
