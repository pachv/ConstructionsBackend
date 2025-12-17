package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/transport/http/responses"
)

const maxAge = 60 * 60 * 24 * 30

func (h *Handler) LoginHandler(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		responses.BadRequest(c, "password or username empty")
		c.Abort()
	}

	user, err := h.authService.LoginUser(username, password)
	if err != nil {
		responses.BadRequest(c, "cant login user : "+err.Error())
		c.Abort()
	}

	sessionId, err := h.sessionService.CreateSession(user.Username, user.Id)
	if err != nil {
		responses.BadRequest(c, "cant create sesstion : "+err.Error())
		c.Abort()
	}

	fmt.Println("setting session " + sessionId)

	c.SetCookie(
		"session_id",
		sessionId,
		maxAge,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) LogoutUser(c *gin.Context) {

	sessionId := c.GetString("sessionId")

	err := h.sessionService.Delete(sessionId)
	if err != nil {
		responses.BadRequest(c, "cant delete session : "+err.Error())
		c.Abort()
		return
	}

	c.SetCookie(
		"session_id", // имя куки
		"",           // пустое значение
		-1,           // MaxAge < 0 → удаление
		"/",          // путь
		"",           // домен (если не задавался, оставляем пустым)
		false,        // secure (как при установке)
		true,         // httpOnly
	)

	c.JSON(200, gin.H{})
}
