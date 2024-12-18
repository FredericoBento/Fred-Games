package handler

import (
	"net/http"

	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/home_views"

	"github.com/a-h/templ"
)

type HomeHandler struct {
	games       []services.GameService
	authService *services.AuthService
}

func NewHomeHandler(gameServices []services.GameService, authService *services.AuthService) *HomeHandler {

	return &HomeHandler{
		games:       gameServices,
		authService: authService,
	}

	// h.navbar = h.getNavbar(false, false)
	// h.isLogged = false
	// h.isAdmin = false

	// return h
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/home":
		h.home(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	}
}

func (h *HomeHandler) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetHome(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type HomeViewProps struct {
	title   string
	content templ.Component
}

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	// apps := h.appManager.GetAppsSortedAlphabetic()
	props := HomeViewProps{
		content: home_views.Home(h.games),
	}
	h.View(w, r, props)
}

func (h *HomeHandler) View(w http.ResponseWriter, r *http.Request, props HomeViewProps) {
	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		var navbar templ.Component
		if IsLogged(r) {
			if IsAdmin(r) {
				navbar = home_views.AdminNavbar()
			} else {
				navbar = home_views.LoggedNavbar()
			}
		} else {
			navbar = home_views.DefaultNavbar()
		}
		views.Page(props.title, navbar, props.content).Render(r.Context(), w)
	}
}

func (h *HomeHandler) getNavbar(isLogged bool, isAdmin bool) models.NavBarStructure {
	startBtns := []models.Button{
		{
			ButtonName: "Games",
			Url:        "/home",
		},
	}

	var endBtns []models.Button
	if isLogged {

		if isAdmin {
			endBtns = []models.Button{
				{
					ButtonName: "Account",
					Childs: []models.Button{
						{
							ButtonName:   "Admin",
							Url:          "/admin/dashboard",
							NotHxRequest: true,
						},
						{
							ButtonName: "Settings",
							Url:        "/settings",
						},
						{
							ButtonName: "Logout",
							Url:        "/logout",
						},
					},
				},
			}

		} else {
			endBtns = []models.Button{
				{
					ButtonName: "Account",
					Childs: []models.Button{
						{
							ButtonName: "Settings",
							Url:        "/settings",
						},
						{
							ButtonName: "Logout",
							Url:        "/logout",
						},
					},
				},
			}
		}
	} else {

		endBtns = []models.Button{
			{
				ButtonName:   "Sign Up",
				Url:          "/sign-up",
				NotHxRequest: true,
			},
			{
				ButtonName:   "Log In",
				Url:          "/sign-in",
				NotHxRequest: true,
			},
		}
	}

	return models.NavBarStructure{
		StartButtons: startBtns,
		EndButtons:   endBtns,
	}

}
