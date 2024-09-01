package handler

import (
	"github.com/FredericoBento/HandGame/internal/views/home"
	"net/http"
)

type HomeHandler struct {
}

func NewHomeHandler() http.Handler {
	return &HomeHandler{}
}

func (hh *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/home":
		hh.home(w, r)
	}
}

func (hh *HomeHandler) home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hh.GetHome(w, r)
	}

}

func (hh *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	props := ViewProps{}
	hh.View(w, r, props)
}

type ViewProps struct {
}

func (hh *HomeHandler) View(w http.ResponseWriter, r *http.Request, props ViewProps) {
	home.Index().Render(r.Context(), w)
}
