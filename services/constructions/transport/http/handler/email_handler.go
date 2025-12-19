package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/transport/http/handler/responses"
)

type SetAdminEmailRequest struct {
	Email string `json:"email"`
}

// POST /admin/email
func (h *Handler) SetAdminEmail(c *gin.Context) {
	var req SetAdminEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(c, "cant bind json: "+err.Error())
		return
	}

	if err := h.adminEmailService.Set(c.Request.Context(), req.Email); err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// GET /admin/email
func (h *Handler) GetAdminEmail(c *gin.Context) {
	email, err := h.adminEmailService.Get(c.Request.Context())
	if err != nil {
		responses.InternalServiceErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"email":  email,
	})
}
