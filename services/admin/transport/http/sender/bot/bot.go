package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AdminData struct {
	AppUrl               string `db:"app_url" json:"appURL"`
	WelcomeText          string `db:"welcome_text" json:"welcomeText"`
	UnknownText          string `db:"unknown_text" json:"unknownText"`
	ReferalActivatedText string `db:"referal_activated_text" json:"referalAtivated"`
	PrayerEngText        string `db:"prayer_eng_text" json:"prayerEngText"`
	PrayerArText         string `db:"prayer_ar_text" json:"prayerArText"`
	ImgRoute             string `db:"bot_img_route" json:"imgRoute"`
	BotImgFilename       string `db:"bot_img_filename" json:"botImgFilename"`
}

func GetBotData() (*AdminData, error) {
	url := "http://bot:8080/get-bot-data"
	const maxRetries = 5

	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		wait := time.Duration(attempt*attempt) * time.Second
		fmt.Printf("Попытка %d не удалась: %v. Повтор через %s...\n", attempt, err, wait)
		time.Sleep(wait)
	}

	if err != nil {
		return nil, fmt.Errorf("все %d попыток неудачны: %w", maxRetries, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var data AdminData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return &data, nil
}
