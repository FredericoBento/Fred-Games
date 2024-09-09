package handler

import (
	"net/http"

	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/home_views"
	"github.com/a-h/templ"
)

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (hh *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hh.Get(w, r)
	}
}

func (hh *HomeHandler) Get(w http.ResponseWriter, r *http.Request) {

	hh.View(w, r, homeViewProps{
		title:   "Sign In",
		content: home_views.Base(),
	})
}

type homeViewProps struct {
	title   string
	content templ.Component
}

func (hh *HomeHandler) View(w http.ResponseWriter, r *http.Request, props homeViewProps) {

	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, props.content).Render(r.Context(), w)
	}
}
