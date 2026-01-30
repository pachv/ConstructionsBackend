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
	adminReviewService      *services.AdminReviewService
	adminDashboardService   *services.AdminDashboardService
	adminGalleryService     *services.AdminGalleryService
	contactsService         *services.ContactsService
	adminSectionService     *services.SiteSectionsAdminService
	adminContactsService    *services.AdminContactsService
	productAdminService     *services.AdminProductService
	adminFontsService       *services.AdminFontsService
}

func New(logger *slog.Logger, userService *services.UserService,
	tokenService *services.TokenService, askQuestionService *services.AskQuestionService,
	callbackService *services.CallbackService, reviewService *services.ReviewService,
	productService *services.ProductService, orderService *services.OrderService,
	certificateService *services.CertificateService, galleryService *services.GalleryService,
	siteSectionService *services.SiteSectionsService, adminEmailService *services.AdminEmailService,
	certificateAdminService *services.CertificatesAdminService, adminReviewService *services.AdminReviewService,
	adminDashboardService *services.AdminDashboardService, adminGalleryService *services.AdminGalleryService,
	contactsService *services.ContactsService, adminSectionService *services.SiteSectionsAdminService,
	adminContactsService *services.AdminContactsService, productAdminService *services.AdminProductService, adminFontsService *services.AdminFontsService) *Handler {
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
		adminReviewService:      adminReviewService,
		adminDashboardService:   adminDashboardService,
		adminGalleryService:     adminGalleryService,
		contactsService:         contactsService,
		adminSectionService:     adminSectionService,
		adminContactsService:    adminContactsService,
		productAdminService:     productAdminService,
		adminFontsService:       adminFontsService,
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

		contacts := apiv1.Group("/contacts")
		{
			contacts.GET("/email", h.GetContactsEmail)
			contacts.GET("/numbers", h.GetContactsNumbers)
			contacts.GET("/addresses", h.GetContactsAddresses)
		}

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

		admin.GET("/dashboard", h.AdminDashboard)

		admin.GET("/reviews", h.AdminGetReviews)
		admin.POST("/reviews", h.AdminCreateReview)
		admin.DELETE("/reviews/:id", h.AdminDeleteReview)
		admin.PUT("/reviews/bulk", h.AdminBulkUpdateReviews)

		galleryAdmin := admin.Group("/gallery")
		{
			galleryAdmin.POST("/categories", h.AdminCreateGalleryCategory)
			galleryAdmin.PUT("/categories/:id", h.AdminUpdateGalleryCategory)
			galleryAdmin.DELETE("/categories/:id", h.AdminDeleteGalleryCategory)

			galleryAdmin.POST("/categories/:id/photos", h.AdminAddGalleryPhoto)
			galleryAdmin.DELETE("/photos/:id", h.AdminDeleteGalleryPhoto)
		}

		adminSections := admin.Group("/sections")
		{
			// ✅ GET
			adminSections.GET("/all", h.AdminGetSectionsAll)
			adminSections.GET("/:slug/full", h.AdminGetSectionFullBySlug)

			// ✅ CREATE & DELETE
			adminSections.POST("/create-form", h.AdminCreateSectionForm)
			adminSections.DELETE("/:id", h.AdminDeleteSection)

			// ✅ UPDATE (multipart)
			// Вариант 1: “нормальный” REST
			adminSections.PUT("/:id", h.AdminUpdateSectionForm)

			// Вариант 2: чтобы совпало с тем, что ты сейчас дергаешь:
			adminSections.PUT("/update/:id", h.AdminUpdateSectionForm)

			// toggles
			adminSections.PATCH("/:id/gallery/toggle", h.AdminToggleSectionGallery)
			adminSections.PATCH("/:id/catalog/toggle", h.AdminToggleSectionCatalog)

			// gallery
			adminSections.POST("/:id/gallery", h.AdminAddSectionGalleryPhoto)
			adminSections.DELETE("/:id/gallery/:photoId", h.AdminDeleteSectionGalleryPhoto)

			// catalog categories
			adminSections.POST("/:id/catalog/categories", h.AdminAddSectionCatalogCategory)
			adminSections.DELETE("/:id/catalog/categories/:categoryId", h.AdminDeleteSectionCatalogCategory)

			// catalog items
			adminSections.POST("/:id/catalog/items", h.AdminAddSectionCatalogItem)
			adminSections.DELETE("/:id/catalog/items/:itemId", h.AdminDeleteSectionCatalogItem)
		}

		contactsAdmin := admin.Group("/contacts")
		{
			contactsAdmin.GET("", h.AdminGetContacts)

			contactsAdmin.PUT("/email", h.AdminSetContactsEmail)

			contactsAdmin.GET("/phones", h.AdminListContactPhones)
			contactsAdmin.PUT("/phones", h.AdminUpsertContactPhone)
			contactsAdmin.DELETE("/phones/:id", h.AdminDeleteContactPhone)
			contactsAdmin.PUT("/phones/reorder", h.AdminReorderContactPhones)

			contactsAdmin.GET("/addresses", h.AdminListContactAddresses)
			contactsAdmin.PUT("/addresses", h.AdminUpsertContactAddress)
			contactsAdmin.DELETE("/addresses/:id", h.AdminDeleteContactAddress)
			contactsAdmin.PUT("/addresses/reorder", h.AdminReorderContactAddresses)
		}

		cat := admin.Group("/catalog")
		{
			// categories
			cat.GET("/categories", h.AdminGetCatalogCategories)
			cat.GET("/categories/by-title/:title", h.AdminGetCatalogCategoryByTitle)
			cat.POST("/categories", h.AdminCreateCatalogCategory)
			cat.PUT("/categories/:id", h.AdminUpdateCatalogCategory)
			cat.DELETE("/categories/:id", h.AdminDeleteCatalogCategory)

			// sections
			cat.GET("/sections", h.AdminGetCatalogSections)
			cat.GET("/sections/by-title/:title", h.AdminGetCatalogSectionByTitle)
			cat.GET("/sections/:id/category", h.AdminGetCatalogSectionCategory)
			cat.POST("/sections", h.AdminCreateCatalogSection)
			cat.PUT("/sections/:id", h.AdminUpdateCatalogSection)
			cat.DELETE("/sections/:id", h.AdminDeleteCatalogSection)

			// products
			cat.GET("/products", h.AdminGetCatalogProducts)
			cat.GET("/products/:id", h.AdminGetCatalogProduct)
			cat.POST("/products", h.AdminCreateCatalogProduct)
			cat.PUT("/products/:id", h.AdminUpdateCatalogProduct)
			cat.DELETE("/products/:id", h.AdminDeleteCatalogProduct)
		}

		fonts := admin.Group("/fonts")
		{
			fonts.GET("", h.AdminListFonts)
			fonts.POST("", h.AdminCreateFont)
			fonts.DELETE("/:id", h.AdminDeleteFont)
			fonts.POST("/:id/select", h.AdminSelectFont)
		}

	}
}
