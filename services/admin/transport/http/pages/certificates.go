package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/certificates"
)

type CertificatesPageData struct {
	Base

	Search  string
	Section string

	Sections []struct {
		Title string
		Slug  string
	}

	Certificates []Cert

	CurrentPage    int
	PageAmount     int
	PagesToDisplay []int
}

type Cert struct {
	ID       string
	Title    string
	FilePath string
}

const baseCertificatesFilePath = "http://localhost:80/api/v1"

// /admin/certificates?page=1&search=...
func (p *Pages) CertificatesPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("./templates/base.html", "./templates/certificates.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	// query params
	search := strings.TrimSpace(c.Query("search"))
	section := strings.TrimSpace(c.Query("section")) // пока просто прокидываем

	page := 1
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}

	// fetch from constructions_service (admin endpoint)
	res, err := certificates.GetAll(c.Request.Context(), page, search)
	if err != nil {
		c.String(http.StatusBadGateway, "failed to fetch certificates: "+err.Error())
		return
	}

	fmt.Println(res)

	// map items
	certificatesItems := make([]Cert, 0, len(res.Items))
	for _, item := range res.Items {
		certificatesItems = append(certificatesItems, Cert{
			ID:       item.ID,
			Title:    item.Title,
			FilePath: baseCertificatesFilePath + item.FilePath,
		})
	}

	// pagination numbers
	ptd := buildPagesToDisplay(res.Page, res.PageAmount)

	data := CertificatesPageData{
		Base:           p.CreateBase(username, "Сертификаты", "certificates"),
		Search:         search,
		Section:        section,
		Sections:       []struct{ Title, Slug string }{}, // подключишь позже
		Certificates:   certificatesItems,
		CurrentPage:    res.Page,
		PageAmount:     res.PageAmount,
		PagesToDisplay: ptd,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// как в users: если страниц много — показываем окно на 10 вокруг текущей
func buildPagesToDisplay(currentPage, pageAmount int) []int {
	if pageAmount <= 1 {
		return []int{1}
	}
	if currentPage < 1 {
		currentPage = 1
	}
	if currentPage > pageAmount {
		currentPage = pageAmount
	}

	const window = 10

	start := currentPage - window/2
	if start < 1 {
		start = 1
	}

	end := start + window - 1
	if end > pageAmount {
		end = pageAmount
		start = end - window + 1
		if start < 1 {
			start = 1
		}
	}

	out := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		out = append(out, i)
	}
	return out
}
