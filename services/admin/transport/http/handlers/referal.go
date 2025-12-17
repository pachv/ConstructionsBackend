package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/sender/referal"
)

func (h *Handler) SetReferal(c *gin.Context) {
	totalAmountStr := c.PostForm("refered_amount")
	purchasedAmountStr := c.PostForm("premium_purchased_amount")
	totalGiftStr := c.PostForm("referral_requests_amount")
	purchasedGiftStr := c.PostForm("premium_days")

	urlReferal := c.PostForm("url")
	englishText := c.PostForm("english_text")
	arabText := c.PostForm("arab_text")

	totalAmount, _ := strconv.ParseFloat(totalAmountStr, 64)
	purchasedAmount, _ := strconv.ParseFloat(purchasedAmountStr, 64)
	totalGift, _ := strconv.ParseFloat(totalGiftStr, 64)
	purchasedGift, _ := strconv.ParseFloat(purchasedGiftStr, 64)

	reqData := &referal.Bonuses{
		TotalAmount:     totalAmount,
		PurchasedAmount: purchasedAmount,
		TotalGift:       totalGift,
		PurchasedGift:   purchasedGift,
		URL:             urlReferal,
		EnglishText:     englishText,
		ArabText:        arabText,
	}

	err := referal.SetReferalData(reqData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка отправки запроса: %v", err)
		return
	}

	// Успех
	c.Status(http.StatusOK)
}
