package handler

import (
	"net/http"

	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/auth_views"
	"github.com/a-h/templ"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/sign-in":
		ah.signIn(w, r)
	case "/sign-up":
		ah.signUp(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	}
}

func (ah *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ah.GetSignIn(w, r)
	}
}

func (ah *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ah.GetSignUp(w, r)
	}
}

func (ah *AuthHandler) GetSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ah.View(w, r, viewProps{
		title:   "Sign Up",
		content: auth_views.SignUpForm(),
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sign Up Get"))
}

func (ah *AuthHandler) GetSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sign In Get"))
}

type viewProps struct {
	title   string
	content templ.Component
}

func (ah *AuthHandler) View(w http.ResponseWriter, r *http.Request, props viewProps) {
	views.Page(props.title, props.content).Render(r.Context(), w)
}
