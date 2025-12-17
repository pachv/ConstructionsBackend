package pages

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
)

type UsersSettingsPageData struct {
	Base

	Search      string
	OrderBy     string
	PageAmount  int
	CurrentPage int

	CreateUserURL string

	AdminUsers     []AdminUser
	PagesToDisplay []int
}

type AdminUser struct {
	Id        string
	Username  string
	DeleteURL string
	EditURL   string
}

func (p *Pages) UsersSettingsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/users_settings.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

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

	users, pageAmount, err := p.authService.GetUsers(page, search, orderBy)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	adminUsers := []AdminUser{}

	editURLBASE := p.Domain + "/admin/settings/users/"
	deleteURLBASE := p.Domain + "/admin-service/admin/delete-admin/"

	for _, u := range users {
		adminUsers = append(adminUsers, AdminUser{
			Id:        u.Id,
			Username:  u.Username,
			EditURL:   editURLBASE + u.Id,
			DeleteURL: deleteURLBASE + u.Id,
		})
	}

	pagesToDisplay := calculatePagesToDisplay(page, pageAmount)

	data := UsersSettingsPageData{
		Base:        p.CreateBase(username, "Settings", "settings"),
		PageAmount:  pageAmount,
		CurrentPage: page,
		AdminUsers:  adminUsers,

		Search:  search,
		OrderBy: orderBy,

		CreateUserURL:  p.Domain + "/admin/settings/users/create",
		PagesToDisplay: pagesToDisplay,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func calculatePagesToDisplay(currentPage, pageAmount int) []int {
	const windowSize = 11
	if pageAmount <= 1 {
		return nil
	}

	start := 1
	end := pageAmount

	if pageAmount > windowSize {
		start = currentPage - windowSize/2
		if start < 1 {
			start = 1
		}
		end = start + windowSize - 1
		if end > pageAmount {
			end = pageAmount
			start = end - windowSize + 1
		}
	}

	pages := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}
