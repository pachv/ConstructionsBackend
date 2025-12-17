package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/prayer"
)

type Prayer struct {
	ID                string `json:"id"`
	PrayerName        string `json:"prayer_name"`
	PrayerLanguageID  string `json:"prayer_language_id"`
	EngName           string `json:"eng_name"`
	PrayerDescription string `json:"prayer_description"`
}

type PrayersData struct {
	Base

	EnglishPrayers []Prayer
	ArabicPrayers  []Prayer
	SetPrayersURL  string
}

const ArabLanguageId = "7b64a96d-1dc9-4cd0-b3f0-59cbfbc9fdf7"

func (p *Pages) PrayersPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/prayers.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	engPrayers := []Prayer{}

	arabPrayers := []Prayer{}

	prayers, err := prayer.GetPrayers()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for _, p := range prayers {

		if p.LanguageId == ArabLanguageId {
			arabPrayers = append(arabPrayers, Prayer{
				ID:                p.ID,
				PrayerName:        p.Name,
				PrayerLanguageID:  p.LanguageId,
				EngName:           p.EngName,
				PrayerDescription: p.Description,
			})
		} else {
			engPrayers = append(engPrayers, Prayer{
				ID:                p.ID,
				PrayerName:        p.Name,
				PrayerLanguageID:  p.LanguageId,
				EngName:           p.EngName,
				PrayerDescription: p.Description,
			})
		}
	}

	data := PrayersData{
		Base:           p.CreateBase(username, "Prayers", " prayers"),
		EnglishPrayers: engPrayers,
		ArabicPrayers:  arabPrayers,
		SetPrayersURL:  p.Domain + "/admin-service/admin/set-prayers",
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
