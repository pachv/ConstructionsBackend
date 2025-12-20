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

type GalleryPhotoDTO struct {
	ID           string `json:"id"`
	CategorySlug string `json:"category_slug"`
	Alt          string `json:"alt"`
	Image        string `json:"image"` // ВАЖНО: это только имя файла (image_path)
	SortOrder    int    `json:"sort_order"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type galleryPhotosResp struct {
	Items []GalleryPhotoDTO `json:"items"`
}

type GalleryCategoryPageData struct {
	Base
	Category GalleryCategoryDTO
	Photos   []GalleryPhotoDTO

	// чтобы в шаблоне строить url картинки через публичный api
	PublicAPI string
}

func (p *Pages) GalleryCategoryPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("./templates/base.html", "./templates/gallery_category.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")
	slug := c.Param("slug")

	// найдём категорию по slug (через список категорий)
	cats, err := p.fetchGalleryCategories(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var cat GalleryCategoryDTO
	found := false
	for _, it := range cats {
		if it.Slug == slug {
			cat = it
			found = true
			break
		}
	}
	if !found {
		c.String(http.StatusNotFound, "category not found")
		return
	}

	photos, err := p.fetchGalleryPhotosBySlug(c, slug)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := GalleryCategoryPageData{
		Base:      p.CreateBase(username, "Галерея — "+cat.Title, "gallery"),
		Category:  cat,
		Photos:    photos,
		PublicAPI: "http://localhost:80", // если у тебя в Base есть PublicAPIBaseURL — замени на него
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func (p *Pages) fetchGalleryPhotosBySlug(c *gin.Context, slug string) ([]GalleryPhotoDTO, error) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, constructionsPublicAPI+"/gallery/"+slug+"/photos", nil)
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
		return nil, fmt.Errorf("gallery photos bad status: %s body=%s", res.Status, string(b))
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// 1) пробуем {"items":[...]}
	var wrap struct {
		Items []GalleryPhotoDTO `json:"items"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil && wrap.Items != nil {
		return wrap.Items, nil
	}

	// 2) пробуем просто [...]
	var arr []GalleryPhotoDTO
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}

	return nil, fmt.Errorf("unexpected gallery photos response: %s", string(raw))
}
