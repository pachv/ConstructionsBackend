package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TonDataRequest struct {
	Wallet string `json:"wallet"`
	URL    string `json:"url"`
}

func (h *Handler) SetTon(c *gin.Context) {
	// Получаем данные из формы
	wallet := c.PostForm("wallet")
	urlStr := c.PostForm("url")

	// Формируем тело запроса
	reqData := TonDataRequest{
		Wallet: wallet,
		URL:    urlStr,
	}

	jsonBytes, err := json.Marshal(reqData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка формирования JSON: %v", err)
		return
	}

	// Отправка POST-запроса на внешний сервис
	resp, err := http.Post("http://payment:8080/admin/ton-data", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка отправки запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.String(http.StatusInternalServerError, "Сервер вернул статус: %s", resp.Status)
		return
	}

	// Успех
	c.Status(http.StatusOK)
}
