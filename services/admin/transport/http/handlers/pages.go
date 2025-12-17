package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) InitPagesHandlers(r *gin.RouterGroup) {
	r.GET("/", h.pages.DashboardPage)
	r.GET("/friends", h.pages.ReferalPage)
	r.GET("/ton", h.pages.TonPage)
	r.GET("/users", h.pages.UsersPage)
	r.GET("/bot", h.pages.BotPage)
	r.GET("/translations", h.pages.TranslationsPage)
	r.GET("/prayers", h.pages.PrayersPage)
	r.GET("/prices", h.pages.PricesPage)

	r.GET("/settings", h.pages.SettingsPage)
	r.GET("/settings/users", h.pages.UsersSettingsPage)
	r.GET("/settings/users/create", h.pages.CreateUserPage)
	r.GET("/settings/users/:id", h.pages.EditUserPage)
	r.GET("/payments", h.pages.PaymentsPage)
}
