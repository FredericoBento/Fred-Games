package handler

import (
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/components"
	"github.com/a-h/templ"
)

type HandGameHandler struct {
	handGameService *services.HandGameService
	log             *slog.Logger
}

func NewHandGameHandler(handGameService *services.HandGameService) *HandGameHandler {
	lo, err := logger.NewHandlerLogger("handgame", "", false)
	if err != nil {
		lo = slog.New(slog.Default().Handler())
		lo.Error(err.Error())
	}

	return &HandGameHandler{
		handGameService: handGameService,
		log:             lo,
	}
}

func (h *HandGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	}
}

func (h *HandGameHandler) Get(w http.ResponseWriter, r *http.Request) {

	// h.View(w, r, HandGameViewProps{
	// content: handgame_views.Home(),
	// })
}

type HandGameViewProps struct {
	title   string
	content templ.Component
}

func (h *HandGameHandler) View(w http.ResponseWriter, r *http.Request, props HandGameViewProps) {

	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, components.DefaultLoggedNavbar(), props.content).Render(r.Context(), w)
	}
}
