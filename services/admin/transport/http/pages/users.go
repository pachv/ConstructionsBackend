package pages

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/user"
)

type UserPageData struct {
	Base

	Search         string
	OrderBy        string
	PageAmount     int
	CurrentPage    int
	Users          []*UserData
	PagesToDisplay []int
}

type UserData struct {
	Id              string
	Username        string
	TelegramId      int64
	FirstLogin      string
	LastAuth        string
	PremiumDaysLeft int64
	RequestsLeft    int64
	Promo           string
	ReferalSenderId string
}

func (p *Pages) UsersPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/users.html",
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

	// Достаём пользователей через сервис
	// users, pageAmount, err := p.adminService.GetUserData(page, search, orderBy)
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	userPagesData, err := user.FetchUsersData(page, search, orderBy)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Подготавливаем отображение
	pagesToDisplay := calculatePagesToDisplay(page, userPagesData.PageAmount)

	var users []*UserData

	for _, u := range userPagesData.Users {
		users = append(users, &UserData{
			Id:              u.Id,
			Username:        u.Username,
			TelegramId:      u.TelegramId,
			FirstLogin:      u.FirstLogin.Format(time.RFC3339),
			LastAuth:        u.LastAuth.Format(time.RFC3339),
			PremiumDaysLeft: u.PremiumDaysLeft,
			RequestsLeft:    u.RequestsLeft,
			Promo:           u.Promo,
			ReferalSenderId: u.ReferalSenderId,
		})
	}

	data := UserPageData{
		Base:        p.CreateBase(username, "Users", "users"),
		PageAmount:  userPagesData.PageAmount,
		CurrentPage: page,
		Users:       users,

		Search:  search,
		OrderBy: orderBy,

		PagesToDisplay: pagesToDisplay,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
