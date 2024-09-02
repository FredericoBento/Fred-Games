package handler

import (
	"net/http"

	"github.com/FredericoBento/HandGame/internal/views/home_views"
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
	props := ViewProps{}
	hh.View(w, r, props)
}

type ViewProps struct {
}

func (hh *HomeHandler) View(w http.ResponseWriter, r *http.Request, props ViewProps) {
	home_views.Index().Render(r.Context(), w)
}
