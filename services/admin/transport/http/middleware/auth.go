package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/internal/repository"
)

type AuthMiddleware struct {
	sessionRepository *repository.SessionRepository
	baseURL           string
}

func NewAuthMiddleware(sessionRepository *repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepository: sessionRepository,
		baseURL:           "http://localhost:8111/admin",
	}
}

func (m *AuthMiddleware) IsSessionCorrectPages(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	userSession, err := m.sessionRepository.GetSessionBySessionId(sessionID)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	if time.Now().After(userSession.ExpiresAt) {
		m.sessionRepository.DeleteSession(sessionID)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	// достается из репозитория сессии есть ли такое, если нет, то редирект на login
	c.Set("username", userSession.UserName)
	c.Set("userId", userSession.UserId)
	c.Set("sessionId", sessionID)

	c.Next()
}

func (m *AuthMiddleware) IsSessionCorrect(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(401, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}

	fmt.Println("sessionId is " + sessionID)

	userSession, err := m.sessionRepository.GetSessionBySessionId(sessionID)
	if err != nil {
		c.JSON(401, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}

	c.Set("user", userSession.UserId)
	c.Next()
}
