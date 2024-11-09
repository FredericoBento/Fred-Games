package handler

import (
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/components"
	"github.com/FredericoBento/HandGame/internal/views/tictactoe_views"
	"github.com/a-h/templ"
)

type TicTacToeHandler struct {
	log *slog.Logger
}

type TicTacToeViewProps struct {
	title   string
	content templ.Component
}

func NewTicTacToeHandler() *TicTacToeHandler {
	lo, err := logger.NewHandlerLogger("TicTacToeHandler", "", false)
	if err != nil {
		lo = slog.New(slog.Default().Handler())
		lo.Error(err.Error())
	}
	return &TicTacToeHandler{
		log: lo,
	}
}

func (h *TicTacToeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/tictactoe/home":
		h.home(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	}
}

func (h *TicTacToeHandler) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getHome(w, r)
	}
}

func (h *TicTacToeHandler) getHome(w http.ResponseWriter, r *http.Request) {
	if !IsLogged(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	h.View(w, r, TicTacToeViewProps{
		content: tictactoe_views.Home(),
	})
}

func (h *TicTacToeHandler) View(w http.ResponseWriter, r *http.Request, props TicTacToeViewProps) {
	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, components.DefaultLoggedNavbar(), props.content).Render(r.Context(), w)
	}
}

// func (h *TicTacToeHandler) getHome(w http.ResponseWriter, r *http.Request) {

// }
