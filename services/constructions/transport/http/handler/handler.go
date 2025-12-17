package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/services"
)

type Handler struct {
	logger *slog.Logger

	userService        *services.UserService
	tokenService       *services.TokenService
	askQuestionService *services.AskQuestionService
}

func New(logger *slog.Logger, userService *services.UserService,
	tokenService *services.TokenService, askQuestionService *services.AskQuestionService) *Handler {
	return &Handler{
		logger:             logger.With("component", "handler"),
		userService:        userService,
		tokenService:       tokenService,
		askQuestionService: askQuestionService,
	}
}

func (h *Handler) InitRoutes(engine *gin.Engine) {

	apiv1 := engine.Group("/api/v1")
	{
		user := apiv1.Group("/user")
		{
			user.POST("/register", h.RegisterUser)
			user.POST("/login", h.Login)
			user.GET("/me", h.AuthMiddleware(), h.Me)
			user.POST("/change-password", h.AuthMiddleware(), h.ChangePassword)
			user.POST("/logout", h.LogOut)
		}

		ratings := apiv1.Group("ratings")
		{
			ratings.GET("")
		}
	}
}
