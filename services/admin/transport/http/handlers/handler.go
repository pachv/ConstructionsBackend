package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/internal/service"
	"github.com/is_backend/services/admin/transport/http/middleware"
	"github.com/is_backend/services/admin/transport/http/pages"
)

type Handler struct {
	logger         *slog.Logger
	authMiddleware *middleware.AuthMiddleware
	pages          *pages.Pages

	// services
	authService    *service.AuthService
	sessionService *service.SessionService
}

func NewHandler(
	logger *slog.Logger,
	authMiddleware *middleware.AuthMiddleware,
	pagesDomain string,
	authService *service.AuthService,
	sessionService *service.SessionService) *Handler {
	return &Handler{
		logger:         logger,
		authMiddleware: authMiddleware,
		pages:          pages.New(pagesDomain, authService),
		authService:    authService,
		sessionService: sessionService,
	}
}

func (h *Handler) Init(e *gin.Engine) {
	admin := e.Group("/admin", h.authMiddleware.IsSessionCorrectPages)
	{
		h.InitPagesHandlers(admin)
	}

	e.GET("/admin/login/", h.pages.LoginPage)

	e.POST("/admin-service/admin/login", h.LoginHandler)

	inside := e.Group("/admin-service/admin", h.authMiddleware.IsSessionCorrect)
	{
		h.InitInsideHandlers(inside)
	}
}
