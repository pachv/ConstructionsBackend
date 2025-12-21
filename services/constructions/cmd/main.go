package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/config"
	"github.com/pachv/constructions/constructions/internal/repositories"
	"github.com/pachv/constructions/constructions/internal/services"
	"github.com/pachv/constructions/constructions/logger"
	"github.com/pachv/constructions/constructions/store"
	"github.com/pachv/constructions/constructions/transport/http/handler"
)

func main() {
	cfg, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewLogger(cfg.Logger.Level)

	logger.Info("starting...")
	logger.Info("log level selected : " + cfg.Logger.Level)

	store := store.NewPostgreSQLStore(cfg, logger)

	store.MustConnect()
	store.MakeMigrations()

	// ! repositories
	userRepository := repositories.NewUserRepository(logger, store.GetDB())
	askQuestionRepository := repositories.NewAskQuestionRepository(logger, store.GetDB())
	callbackRepository := repositories.NewCallbackRepository(logger, store.GetDB())
	reviewRepository := repositories.NewReviewRepository(logger, store.GetDB(), cfg.Domain)
	productRepository := repositories.NewProductRepository(store.GetDB(), logger)
	orderRepository := repositories.NewOrderRepository(store.GetDB(), logger)
	certificatesRepository := repositories.NewCertificateRepository(store.GetDB(), logger)
	galleryRepository := repositories.NewGalleryRepository(store.GetDB(), logger, cfg.Domain)
	sectionsRepository := repositories.NewSiteSectionsRepository(store.GetDB(), cfg.Domain)
	emailRepository := repositories.NewAdminEmailRepository(store.GetDB())
	contactsRepository := repositories.NewContactsRepository(store.GetDB())

	EMAIL_TO_SEND_MAIL := emailRepository.GetEmail()

	fmt.Println(" email is  " + EMAIL_TO_SEND_MAIL)

	// EMAIL_TO_SEND_MAIL := cfg.Email.NotifyEmail

	// ! services
	passwordService := services.NewPasswordService(10)
	userService := services.NewUserService(userRepository, logger, passwordService)
	tokenService := services.NewTokenService(cfg.JWT.Secret)
	mailSendingService := services.NewMailSendingService(cfg.Email.From, cfg.Email.Host, cfg.Email.Password, cfg.Email.Port)
	askQuestionService := services.NewAskQuestionService(
		logger,
		askQuestionRepository,
		mailSendingService,
		[]string{EMAIL_TO_SEND_MAIL},
		"./templates/ask_question.html")
	callbackService := services.NewCallbackService(
		logger,
		callbackRepository,
		mailSendingService,
		cfg.Email.NotifyEmail,
		"templates/callback.html",
	)
	reviewService := services.NewReviewService(logger, reviewRepository)
	productService := services.NewProductService(productRepository)
	orderService := services.NewOrderService(orderRepository, mailSendingService, logger, []string{EMAIL_TO_SEND_MAIL})
	certificatesService := services.NewCertificateService(certificatesRepository, logger)
	galleryService := services.NewGalleryService(galleryRepository)
	sectionsService := services.NewSiteSectionsService(sectionsRepository)
	emailService := services.NewAdminEmailService(emailRepository)
	adminService := services.NewCertificatesAdminService(store.GetDB())

	reviewAdminService := services.NewAdminReviewService(store.GetDB(), cfg.Domain)
	adminDashboardService := services.NewAdminDashboardService(store.GetDB())
	adminGalleryService := services.NewAdminGalleryService(store.GetDB(), "./uploads/gallery")
	contactsService := services.NewContactsService(contactsRepository)
	adminSectionService := services.NewSiteSectionsAdminService(store.GetDB(), cfg.Domain)

	// ! handler

	eng := gin.Default()

	handler := handler.New(logger, userService,
		tokenService, askQuestionService,
		callbackService, reviewService, productService,
		orderService, certificatesService, galleryService,
		sectionsService, emailService, adminService, reviewAdminService,
		adminDashboardService, adminGalleryService,
		contactsService, adminSectionService)
	handler.InitRoutes(eng)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: eng,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("user forcefully shutdown")
		logger.Error("constructions forcefully shutdown")
	}

	logger.Info("constructions shutdown")

}
