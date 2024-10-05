package handler

import (
	"log/slog"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/home_views"
	"github.com/a-h/templ"
)

type HandGameHandler struct {
	navbar models.NavBarStructure
	log    *slog.Logger
}

func NewHandGameHandler() *HandGameHandler {
	lo, err := logger.NewHandlerLogger("handgame", "", false)
	if err != nil {
		lo = slog.New(slog.Default().Handler())
		lo.Error(err.Error())
	}

	h := &HandGameHandler{
		log: lo,
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

	h.View(w, r, handGameViewProps{
		title:       "Sign In",
		headerTitle: "Fred's Apps",
		content:     home_views.Base(),
	})
}

type handGameViewProps struct {
	title       string
	headerTitle string
	content     templ.Component
}

func (h *HandGameHandler) View(w http.ResponseWriter, r *http.Request, props handGameViewProps) {

	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, props.headerTitle, h.navbar, props.content).Render(r.Context(), w)
	}
}
