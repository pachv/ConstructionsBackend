package pages

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/bot"
)

type BotData struct {
	Base

	WelcomeText          string
	UnknownText          string
	ReferalActivatedText string

	EnglishPrayerText string
	ArabPrayerText    string

	ImageURL string

	AppURL string

	UpdateTextURL  string
	UpdateImageURL string
}

func (p *Pages) BotPage(c *gin.Context) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/bot.html",
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := c.GetString("username")

	botData, err := bot.GetBotData()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := BotData{
		Base:                 p.CreateBase(username, "Bot settings", " bot"),
		ImageURL:             botData.AppUrl + botData.ImgRoute + botData.BotImgFilename,
		AppURL:               botData.AppUrl,
		WelcomeText:          botData.WelcomeText,
		UnknownText:          botData.UnknownText,
		EnglishPrayerText:    botData.PrayerEngText,
		ArabPrayerText:       botData.PrayerArText,
		ReferalActivatedText: botData.ReferalActivatedText,

		UpdateTextURL:  p.Domain + "/admin-service/admin/update-bot-data",
		UpdateImageURL: p.Domain + "/admin-service/admin/upload-bot-img",
	}
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
