package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type ProductMock struct {
	ID          string
	Title       string
	PriceRub    int
	ImageURL    string
	Description string
}

type ProductsPageData struct {
	Base
	Items []ProductMock
}

func (p *Pages) ProductsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("./templates/base.html", "./templates/products.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	// ✅ Мок-товары (инструменты) + фото с сайта Unsplash (динамические)
	items := []ProductMock{
		{
			ID:          "prd-1",
			Title:       "Молоток столярный 450г",
			PriceRub:    690,
			ImageURL:    "http://localhost:80/api/v1/sections/picture/build-1.jpg",
			Description: "Удобная рукоять, баланс, для дома и стройки.",
		},
		{
			ID:          "prd-2",
			Title:       "Шуруповёрт аккумуляторный 18V",
			PriceRub:    4990,
			ImageURL:    "http://localhost:80/api/v1/sections/picture/build-1.jpg",
			Description: "2 скорости, подсветка, кейс.",
		},
		{
			ID:          "prd-3",
			Title:       "Набор отвёрток 12 шт",
			PriceRub:    1290,
			ImageURL:    "http://localhost:80/api/v1/sections/picture/build-1.jpg",
			Description: "Крест/шлиц, магнитные наконечники.",
		},
		{
			ID:          "prd-4",
			Title:       "Рулетка 5м",
			PriceRub:    390,
			ImageURL:    "http://localhost:80/api/v1/sections/picture/build-1.jpg",
			Description: "Стопор, клипса, ударопрочный корпус.",
		},
		{
			ID:          "prd-5",
			Title:       "Пассатижи 180мм",
			PriceRub:    790,
			ImageURL:    "http://localhost:80/api/v1/sections/picture/build-1.jpg",
			Description: "Закалённая сталь, мягкие ручки.",
		},
		{
			ID:          "prd-6",
			Title:       "Углошлифовальная машина 125мм",
			PriceRub:    3990,
			ImageURL:    "http://localhost:80/api/v1/sections/picture/build-1.jpg",
			Description: "Защита кожуха, удобный хват.",
		},
	}

	data := ProductsPageData{
		Base:  p.CreateBase(username, "Товары", "products"),
		Items: items,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
