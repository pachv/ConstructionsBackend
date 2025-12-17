package pages

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	paymentsdata "github.com/is_backend/services/admin/transport/http/sender/payments-data"
)

type PaymentPageData struct {
	Base

	Search         string
	OrderBy        string
	PageAmount     int
	CurrentPage    int
	Payments       []*PaymentData
	PagesToDisplay []int
}

type PaymentData struct {
	Id            string
	UserId        string
	PriceType     string
	PriceAmount   string
	RevardType    string
	RevardAmount  string
	PaymentStatus string
	PaymentTime   string
}

func (p *Pages) PaymentsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/payments.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Получаем query параметры
	pageString := c.Query("page")
	search := c.Query("search")
	orderBy := c.Query("orderBy")

	var page int
	if pageString != "" {
		page, err = strconv.Atoi(pageString)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		page = 1
	}

	username := c.GetString("username")

	paymentPagesData, err := paymentsdata.FetchPaymentsData(page, search, orderBy)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	pagesToDisplay := calculatePagesToDisplay(page, paymentPagesData.PageAmount)

	var payments []*PaymentData

	for _, pmt := range paymentPagesData.Payments {
		payments = append(payments, &PaymentData{
			Id:            pmt.Id,
			UserId:        pmt.UserId,
			PriceType:     pmt.PriceType,
			PriceAmount:   pmt.PriceAmount,
			RevardType:    pmt.RevardType,
			RevardAmount:  pmt.RevardAmount,
			PaymentStatus: pmt.PaymentStatus,
			PaymentTime:   pmt.PaymentTime,
		})
	}

	data := PaymentPageData{
		Base:        p.CreateBase(username, "Payments", "payments"),
		PageAmount:  paymentPagesData.PageAmount,
		CurrentPage: page,
		Payments:    payments,

		Search:  search,
		OrderBy: orderBy,

		PagesToDisplay: pagesToDisplay,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
