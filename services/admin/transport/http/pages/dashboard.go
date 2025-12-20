package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender"
)

type DashboardData struct {
	Base

	TotalUsers      int
	TotalOrders     int
	OrdersToday     int
	OrdersLastMonth int
}

func (p *Pages) DashboardPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/dashboard.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stats, err := sender.GetConstructionsDashboardStats(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	data := DashboardData{
		Base:            p.CreateBase(username, "Главная", " dashboard"),
		TotalUsers:      stats.TotalUsers,
		TotalOrders:     stats.TotalOrders,
		OrdersToday:     stats.OrdersToday,
		OrdersLastMonth: stats.OrdersLastMonth,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
