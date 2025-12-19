package pages

import (
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

type ReviewMock struct {
	ID         string
	Name       string
	Position   string
	Text       string
	Rating     int
	CanPublish bool
	CreatedAt  time.Time
}

type ReviewsMockPageData struct {
	Base

	Items []ReviewMock
}

func (p *Pages) ReviewsMockPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/reviews.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	// ✅ тестовые данные (без сервисов)
	items := []ReviewMock{
		{
			ID:         "rev-1",
			Name:       "Александр",
			Position:   "Заказчик",
			Text:       "Сделали всё быстро и аккуратно. Металлоконструкции пришли в срок, качество огонь.",
			Rating:     5,
			CanPublish: true,
			CreatedAt:  time.Now().Add(-48 * time.Hour),
		},
		{
			ID:         "rev-2",
			Name:       "Марина",
			Position:   "Дизайнер",
			Text:       "Хороший сервис, но хотелось бы чуть быстрее по ответам. В целом рекомендую.",
			Rating:     4,
			CanPublish: true,
			CreatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         "rev-3",
			Name:       "Илья",
			Position:   "",
			Text:       "Цена норм, работа норм. Был небольшой косяк по упаковке, но исправили.",
			Rating:     4,
			CanPublish: false,
			CreatedAt:  time.Now().Add(-8 * time.Hour),
		},
		{
			ID:         "rev-4",
			Name:       "Ольга",
			Position:   "Покупатель",
			Text:       "Не понравилось: задержка доставки на 2 дня. Качество ок.",
			Rating:     3,
			CanPublish: false,
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
	}

	data := ReviewsMockPageData{
		Base:  p.CreateBase(username, "Отзывы", "reviews"),
		Items: items,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
