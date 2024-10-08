package handler

import (
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/app/pong"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/pong_views"
	"github.com/a-h/templ"
)

type PongHandler struct {
	pongApp     *pong.PongApp
	navbar      models.NavBarStructure
	log         *slog.Logger
	authService *services.AuthService
}

func NewPongHandler(pongApp *pong.PongApp, authService *services.AuthService) *PongHandler {
	lo, err := logger.NewHandlerLogger("pong", "", false)
	if err != nil {
		lo = slog.New(slog.Default().Handler())
		lo.Error(err.Error())
	}

	h := &PongHandler{
		pongApp:     pongApp,
		log:         lo,
		authService: authService,
	}

	h.setupNavbar()

	return h
}

func (h *PongHandler) setupNavbar() {
	startBtns := []models.Button{
		{
			ButtonName: "Home",
			Url:        "/home",
		},
	}

	endBtns := []models.Button{
		{
			ButtonName: "Account",
			Childs: []models.Button{
				{
					ButtonName: "Logout",
					Url:        "/logout",
				},
			},
		},
	}

	navbar := models.NavBarStructure{
		StartButtons: startBtns,
		EndButtons:   endBtns,
	}

	h.navbar = navbar
}

func (h *PongHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/pong/home":
		h.home(w, r)
	case "/pong/join-game":
		h.joinGame(w, r)
	case "/pong/create-game":
		h.createGame(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	}
}

func (h *PongHandler) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getHome(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func (h *PongHandler) createGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.postCreateGame(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}
func (h *PongHandler) joinGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getJoinGame(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func (h *PongHandler) getHome(w http.ResponseWriter, r *http.Request) {
	h.View(w, r, PongViewProps{
		content: pong_views.Home(),
	})
}

func (h *PongHandler) postCreateGame(w http.ResponseWriter, r *http.Request) {
	u, isLogged := h.authService.IsLogged(r)
	if !isLogged {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	h.pongApp.CreateGame(u)

	h.View(w, r, PongViewProps{
		content: pong_views.Home(),
	})
}

func (h *PongHandler) getJoinGame(w http.ResponseWriter, r *http.Request) {
	u, isLogged := h.authService.IsLogged(r)
	if !isLogged {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	h.View(w, r, PongViewProps{
		content: pong_views.Home(),
	})
}

type PongViewProps struct {
	title   string
	content templ.Component
}

func (h *PongHandler) View(w http.ResponseWriter, r *http.Request, props PongViewProps) {

	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, h.navbar, props.content).Render(r.Context(), w)
	}
}
