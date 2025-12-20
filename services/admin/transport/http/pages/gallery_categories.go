package pages

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

const constructionsPublicAPI = "http://constructions_service:8080/api/v1"

type GalleryCategoryDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at,omitempty"`
}

type galleryCategoriesResp struct {
	Items []GalleryCategoryDTO `json:"items"`
}

type GalleryPageData struct {
	Base
	Items []GalleryCategoryDTO
}

func (p *Pages) GalleryPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("./templates/base.html", "./templates/gallery.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	items, err := p.fetchGalleryCategories(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := GalleryPageData{
		Base:  p.CreateBase(username, "Галерея", "gallery"),
		Items: items,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func (p *Pages) fetchGalleryCategories(c *gin.Context) ([]GalleryCategoryDTO, error) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, constructionsPublicAPI+"/gallery/categories", nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 10 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("gallery categories bad status: %s body=%s", res.Status, string(b))
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// 1) пробуем {"items":[...]}
	var wrap struct {
		Items []GalleryCategoryDTO `json:"items"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil && wrap.Items != nil {
		return wrap.Items, nil
	}

	// 2) пробуем просто [...]
	var arr []GalleryCategoryDTO
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}

	return nil, fmt.Errorf("unexpected gallery categories response: %s", string(raw))
}
