package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/transport/http/handler/responses"
)

type RegisterUserRequest struct {
	Surname     string `json:"surname"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Fathername  string `json:"fathername"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

const (
	refreshTokenMaxAge = 60 * 60 * 24 * 30
)

func (h *Handler) RegisterUser(c *gin.Context) {
	fmt.Println("register user start")

	var req RegisterUserRequest

	if err := c.BindJSON(&req); err != nil {
		h.logger.Error("Cant bind register user : " + err.Error())
		responses.BadRequestResponse(c, "Cant bind register user  ")
		c.Abort()
		return
	}

	if req.Surname == "" || req.Name == "" || req.Login == "" || req.Fathername == "" || req.Email == "" || req.PhoneNumber == "" || req.Password == "" {
		responses.BadRequestResponse(c, "not enough data")
		h.logger.Error("not enough data")
		c.Abort()
		return
	}

	fmt.Println("after heck")

	userId, err := h.userService.RegisterUser(req.Surname, req.Name, req.Login, req.Fathername, req.Email, req.PhoneNumber, req.Password)
	if err != nil {
		responses.InternalServiceErrorResponse(c, err.Error())
		h.logger.Error("not enough data")
		c.Abort()
		return
	}

	fmt.Println("after service")

	token, err := h.tokenService.CreateRefreshToken(userId, req.Login)
	if err != nil {
		responses.InternalServiceErrorResponse(c, err.Error())
		h.logger.Error("not enough data")
		c.Abort()
		return
	}

	fmt.Println("after token")

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		Path:     "/",
		MaxAge:   refreshTokenMaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)

	fmt.Println("end")

	responses.OkResponse(c, gin.H{
		"status": "ok",
	})

}

type AuthUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	var req AuthUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	u, err := h.userService.Login(req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	refreshToken, err := h.tokenService.CreateRefreshToken(u.Id, u.LoginName)
	if err != nil {
		h.logger.Error("create refresh token error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.SetCookie(
		"refreshToken",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

func (s *Handler) LogOut(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)

}

func (h *Handler) Me(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, _ := userIDAny.(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u, err := h.userService.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, _ := userIDAny.(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.userService.ChangePassword(userID, req.Password); err != nil {
		h.logger.Error("change password error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
