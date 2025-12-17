package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/is_backend/services/admin/internal/config"
	"github.com/is_backend/services/admin/internal/logger"
	"github.com/is_backend/services/admin/internal/repository"
	"github.com/is_backend/services/admin/internal/service"
	"github.com/is_backend/services/admin/store"
	"github.com/is_backend/services/admin/transport/http/handlers"
	"github.com/is_backend/services/admin/transport/http/middleware"
)

func main() {

	l := logger.NewSlogLogger("DEBUG")
	cfg := config.MustLoadConfig("./.env")

	store := store.NewPostgreSQLStore(cfg)
	store.Connect()
	store.MakeMigrations()

	// ! repositories
	sessionReposotory := repository.NewSessionRepository(store.GetDB())
	authRepo := repository.NewUserRepository(store.GetDB())

	// ! middlewares
	authMiddleware := middleware.NewAuthMiddleware(sessionReposotory)

	// ! services
	passwordHashedService := service.NewPasswordHasherService()
	authService := service.NewAuthService(authRepo, passwordHashedService, l)
	sessionService := service.NewSessionService(sessionReposotory, l)

	h := handlers.NewHandler(l, authMiddleware, cfg.Server.Domain, authService, sessionService)

	// gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.LoadHTMLGlob("./templates/*")

	h.Init(g)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: g,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	l.Info("server started successfully on port " + cfg.Server.Port)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit

}
