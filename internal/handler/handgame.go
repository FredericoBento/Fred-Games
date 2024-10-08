package handler

import (
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/app/handgame"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/handgame_views"
	"github.com/a-h/templ"
)

type HandGameHandler struct {
	handGameApp *handgame.HandGameApp
	navbar      models.NavBarStructure
	log         *slog.Logger
}

func NewHandGameHandler(hgApp *handgame.HandGameApp) *HandGameHandler {
	lo, err := logger.NewHandlerLogger("handgame", "", false)
	if err != nil {
		lo = slog.New(slog.Default().Handler())
		lo.Error(err.Error())
	}

	h := &HandGameHandler{
		handGameApp: hgApp,
		log:         lo,
	}

	h.setupNavbar()

	return h
}

func (h *HandGameHandler) setupNavbar() {
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

func (h *HandGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	}
}

func (h *HandGameHandler) Get(w http.ResponseWriter, r *http.Request) {

	h.View(w, r, HandGameViewProps{
		content: handgame_views.Home(h.handGameApp.GetRooms()),
	})
}

type HandGameViewProps struct {
	title   string
	content templ.Component
}

func (h *HandGameHandler) View(w http.ResponseWriter, r *http.Request, props HandGameViewProps) {

	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, h.navbar, props.content).Render(r.Context(), w)
	}
}
