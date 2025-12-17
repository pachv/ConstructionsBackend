package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/translations"
)

type TranslationsData struct {
	Base

	JSONData           template.JS
	SetTranslationsURL string
}

func (p *Pages) TranslationsPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/translations.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	translations, err := translations.GetTranslations()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	var js interface{}
	if err := json.Unmarshal([]byte(translations), &js); err != nil {
		fmt.Println("не валидный json : " + err.Error())
		c.String(http.StatusInternalServerError, "Invalid JSON: "+err.Error())
		return
	}

	fmt.Println("")

	jsonData, err := json.Marshal(js) // сериализуем обратно для JS
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := TranslationsData{
		Base:               p.CreateBase(username, "Translations", " translations"),
		JSONData:           template.JS(jsonData),
		SetTranslationsURL: p.Domain + "/admin-service/admin/set-translationss",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
