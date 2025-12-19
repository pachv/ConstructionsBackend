package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type GalleryItemMock struct {
	ID       string
	Title    string
	ImageURL string
}

type GalleryPageData struct {
	Base
	Items []GalleryItemMock
}

func (p *Pages) GalleryPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("./templates/base.html", "./templates/gallery.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	items := []GalleryItemMock{
		{ID: "gal-1", Title: "Цех — снаружи", ImageURL: "http://localhost:80/api/v1/sections/picture/build-1.jpg"},
		{ID: "gal-2", Title: "Производство", ImageURL: "http://localhost:80/api/v1/sections/picture/build-1.jpg"},
		{ID: "gal-3", Title: "Сварка", ImageURL: "http://localhost:80/api/v1/sections/picture/build-1.jpg"},
		{ID: "gal-4", Title: "Инструменты", ImageURL: "http://localhost:80/api/v1/sections/picture/build-1.jpg"},
	}

	data := GalleryPageData{
		Base:  p.CreateBase(username, "Галерея", "gallery"),
		Items: items,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
