package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/pachv/constructions/constructions/internal/services"
)

type Handler struct {
	logger *slog.Logger

	userService             *services.UserService
	tokenService            *services.TokenService
	askQuestionService      *services.AskQuestionService
	callbackService         *services.CallbackService
	reviewService           *services.ReviewService
	productService          *services.ProductService
	orderService            *services.OrderService
	certificateService      *services.CertificateService
	galleryService          *services.GalleryService
	siteSectionService      *services.SiteSectionsService
	adminEmailService       *services.AdminEmailService
	certificateAdminService *services.CertificatesAdminService
	reviewAdminService      *services.AdminReviewService
}

func New(logger *slog.Logger, userService *services.UserService,
	tokenService *services.TokenService, askQuestionService *services.AskQuestionService,
	callbackService *services.CallbackService, reviewService *services.ReviewService,
	productService *services.ProductService, orderService *services.OrderService,
	certificateService *services.CertificateService, galleryService *services.GalleryService,
	siteSectionService *services.SiteSectionsService, adminEmailService *services.AdminEmailService,
	certificateAdminService *services.CertificatesAdminService, reviewAdminService *services.AdminReviewService) *Handler {
	return &Handler{
		logger:                  logger.With("component", "handler"),
		userService:             userService,
		tokenService:            tokenService,
		askQuestionService:      askQuestionService,
		callbackService:         callbackService,
		reviewService:           reviewService,
		productService:          productService,
		orderService:            orderService,
		certificateService:      certificateService,
		galleryService:          galleryService,
		siteSectionService:      siteSectionService,
		adminEmailService:       adminEmailService,
		certificateAdminService: certificateAdminService,
		reviewAdminService:      reviewAdminService,
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

		apiv1.GET("/sections", h.GetSectionsAll)
		apiv1.GET("/sections/:slug", h.GetSectionBySlug)

		apiv1.GET("/sections/picture/:name", h.GetSectionMainPicture)
		apiv1.GET("/sections/gallery/picture/:name", h.GetSectionGalleryPicture)
		apiv1.GET("/catalog/picture/:name", h.GetCatalogPicture)

	}

	admin := engine.Group("/admin")
	{
		admin.GET("/email", h.GetAdminEmail)
		admin.POST("/email", h.SetAdminEmail)

		admin.GET("/certificates", h.AdminGetAllCertificates)
		admin.POST("/certificates", h.AdminCreateCertificate)
		admin.PUT("/certificates/:id", h.AdminUpdateCertificate)
		admin.DELETE("/certificates/:id", h.AdminDeleteCertificate)

		admin.GET("/certificates/file/:name", h.AdminGetCertificateFile)

		// admin.GET("/reviews", h.AdminGetAllReviews)
		// admin.POST("/reviews/update", h.AdminUpdateReview)
		// admin.POST("/reviews/delete/:id", h.AdminDeleteReview)
	}
}
