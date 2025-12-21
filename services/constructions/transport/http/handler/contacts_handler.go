package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /contacts/email
func (h *Handler) GetContactsEmail(c *gin.Context) {
	out, err := h.contactsService.GetEmail(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contacts email"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// GET /contacts/numbers
func (h *Handler) GetContactsNumbers(c *gin.Context) {
	out, err := h.contactsService.GetNumbers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contacts numbers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

// GET /contacts/addresses
func (h *Handler) GetContactsAddresses(c *gin.Context) {
	out, err := h.contactsService.GetAddresses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contacts addresses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}
