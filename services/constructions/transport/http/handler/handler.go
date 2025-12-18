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
	callbackService    *services.CallbackService
	reviewService      *services.ReviewService
	productService     *services.ProductService
	orderService       *services.OrderService
	certificateService *services.CertificateService
	galleryService     *services.GalleryService
}

func New(logger *slog.Logger, userService *services.UserService,
	tokenService *services.TokenService, askQuestionService *services.AskQuestionService,
	callbackService *services.CallbackService, reviewService *services.ReviewService,
	productService *services.ProductService, orderService *services.OrderService,
	certificateService *services.CertificateService, galleryService *services.GalleryService) *Handler {
	return &Handler{
		logger:             logger.With("component", "handler"),
		userService:        userService,
		tokenService:       tokenService,
		askQuestionService: askQuestionService,
		callbackService:    callbackService,
		reviewService:      reviewService,
		productService:     productService,
		orderService:       orderService,
		certificateService: certificateService,
		galleryService:     galleryService,
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

		ratings := apiv1.Group("/ratings")
		{
			ratings.GET("/add")
		}

		email := apiv1.Group("/email")
		{
			email.POST("/ask-question", h.AskQuestion)
			email.POST("/callback", h.Callback)
			email.POST("/create-order", h.CreateOrder)
		}

		reviews := apiv1.Group("/reviews")
		{
			reviews.POST("", h.CreateReview)
			reviews.GET("", h.GetPublishedReviews)
			reviews.GET("/picture/:name", h.GetReviewPicture)
		}

		products := apiv1.Group("/products")
		{
			products.GET("/categories", h.GetAllCategories)
			products.GET("/sections", h.GetAllSections)
			products.GET("", h.GetAllProducts)
			products.GET("/picture/:image", h.GetProductPicture)
		}

		certs := apiv1.Group("/certificates")
		{
			certs.GET("", h.GetAllCertificates)            // список: [{title, file_path}]
			certs.GET("/file/:name", h.GetCertificateFile) // отдать файл
		}

		// api/v1
		gallery := apiv1.Group("/gallery")
		{
			gallery.GET("/categories", h.GetGalleryCategories)
			gallery.GET("/:slug/photos", h.GetGalleryPhotosByCategory)
			gallery.GET("/picture/:image", h.GetGalleryPicture)
		}

	}
}
