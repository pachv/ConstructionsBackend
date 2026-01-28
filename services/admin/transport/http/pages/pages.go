package pages

import "github.com/is_backend/services/admin/internal/service"

const (
	TEMPLATE_FOLDER_PATH = "./templates"
	PUBLIC_API_BASE      = "http://localhost:80"
)

type Pages struct {
	authService *service.AuthService

	templatesFolderPath string
	Domain              string

	UserLogoutURL string

	DashboardURL    string
	UsersURL        string
	FriendsURL      string
	PricesURL       string
	BotURL          string
	PaymentsURL     string
	TonURL          string
	PrayersURL      string
	TranslationsURL string
	SettingsURL     string

	ProductsURL     string
	GalleryURL      string
	CertificatesURL string

	PublicAPIBaseURL string
}

type Base struct {
	Title  string
	Active string

	FaviconURL string
	LogoURL    string

	Username      string
	UserLogoutURL string

	DashboardURL    string
	UsersURL        string
	FriendsURL      string
	PricesURL       string
	BotURL          string
	PaymentsURL     string
	TonURL          string
	PrayersURL      string
	TranslationsURL string
	SettingsURL     string
	ReviewsURL      string
	EmailURL        string
	FontsURL        string

	SectionsURL string
	ContactsURL string

	ProductsURL     string
	GalleryURL      string
	CertificatesURL string
}

func (p *Pages) CreateBase(username, title, active string) Base {
	return Base{
		Username:        username,
		Title:           title,
		Active:          active,
		DashboardURL:    p.Domain + "/admin",
		FriendsURL:      p.Domain + "/admin/friends",
		TonURL:          p.Domain + "/admin/ton",
		UsersURL:        p.Domain + "/admin/users?page=1",
		BotURL:          p.Domain + "/admin/bot",
		TranslationsURL: p.Domain + "/admin/translations",
		SettingsURL:     p.Domain + "/admin/settings",
		PrayersURL:      p.Domain + "/admin/prayers",
		PricesURL:       p.Domain + "/admin/prices",
		PaymentsURL:     p.Domain + "/admin/payments",
		ReviewsURL:      p.Domain + "/admin/reviews",
		EmailURL:        p.Domain + "/admin/email",

		ProductsURL:     p.Domain + "/admin/products",
		GalleryURL:      p.Domain + "/admin/gallery",
		CertificatesURL: p.Domain + "/admin/certificates?page=1",
		SectionsURL:     p.Domain + "/admin/sections",
		ContactsURL:     p.Domain + "/admin/contacts",
		FontsURL:        p.Domain + "/admin/fonts",

		FaviconURL:    p.Domain + "/admin-service/admin/favicon",
		LogoURL:       p.Domain + "/admin-service/admin/logo",
		UserLogoutURL: p.Domain + "/admin-service/admin/logout",
	}
}

func New(Domain string, authService *service.AuthService) *Pages {
	return &Pages{

		authService:         authService,
		templatesFolderPath: TEMPLATE_FOLDER_PATH,
		Domain:              Domain,
		PublicAPIBaseURL:    PUBLIC_API_BASE,
	}
}
