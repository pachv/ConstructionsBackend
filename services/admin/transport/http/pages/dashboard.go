package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/dashboard"
	"github.com/is_backend/services/admin/transport/http/sender/user"
)

type DashboardData struct {
	Base

	TotalUsers          int
	ActiveSubscriptions int
	RequestsSold        int
	PremiumSold         int
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

	dashboardUserData, err := user.GetDashboardUserData()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	dashboardPaymentData, err := dashboard.GetDashboardPayments()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	data := DashboardData{
		Base:                p.CreateBase(username, "Dashboard", " dashboard"),
		TotalUsers:          dashboardUserData.TotalUsers,
		ActiveSubscriptions: dashboardUserData.ActiveSubscriptionAmount,
		RequestsSold:        dashboardPaymentData.PremiumSold,
		PremiumSold:         dashboardPaymentData.RequestsSold,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
