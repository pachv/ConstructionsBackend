package pages

import "github.com/is_backend/services/admin/internal/service"

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

		FaviconURL:    p.Domain + "/admin-service/admin/favicon",
		LogoURL:       p.Domain + "/admin-service/admin/logo",
		UserLogoutURL: p.Domain + "/admin-service/admin/logout",
	}
}

func New(Domain string, authService *service.AuthService) *Pages {
	return &Pages{

		authService:         authService,
		templatesFolderPath: "./templates",
		Domain:              Domain,
	}
}
