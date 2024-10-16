package handler

import (
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/services/pong"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/components"
	"github.com/FredericoBento/HandGame/internal/views/pong_views"
	"github.com/a-h/templ"
)

type PongHandler struct {
	pongService *pong.PongService
	log         *slog.Logger
}

func NewPongHandler(pongService *pong.PongService) *PongHandler {
	lo, err := logger.NewHandlerLogger("PongHandler", "", false)
	if err != nil {
		lo = slog.New(slog.Default().Handler())
		lo.Error(err.Error())
	}

	return &PongHandler{
		pongService: pongService,
		log:         lo,
	}
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
	if !IsLogged(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	h.View(w, r, PongViewProps{
		content: pong_views.Home(),
	})
}

func (h *PongHandler) postCreateGame(w http.ResponseWriter, r *http.Request) {
	if !IsLogged(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	h.View(w, r, PongViewProps{
		content: pong_views.Home(),
	})
}

func (h *PongHandler) getJoinGame(w http.ResponseWriter, r *http.Request) {
	if !IsLogged(r) {
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
		views.Page(props.title, components.DefaultLoggedNavbar(), props.content).Render(r.Context(), w)
	}
}
