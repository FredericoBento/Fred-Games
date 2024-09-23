package handler

import (
	"context"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/auth_views"
	"github.com/a-h/templ"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
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
	case http.MethodPost:
		ah.PostSignUp(w, r)
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
}

func (ah *AuthHandler) GetSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ah.View(w, r, viewProps{
		title:   "Sign In",
		content: auth_views.SignInForm(),
	})
}

func (ah *AuthHandler) PostSignUp(w http.ResponseWriter, r *http.Request) {
	type SignUpForm struct {
		username         string
		password         string
		repeatedPassword string
	}

	data := SignUpForm{}
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")
	data.repeatedPassword = r.FormValue("repeat_password")

	if data.username == "" || data.password == "" || data.repeatedPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := ah.userService.CreateUser(context.TODO(), &models.User{ID: 0, Username: data.username, Password: data.password})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
	return

}

type viewProps struct {
	title   string
	content templ.Component
}

func (ah *AuthHandler) View(w http.ResponseWriter, r *http.Request, props viewProps) {
	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		var aux map[string]string
		views.Page(props.title, "", aux, props.content).Render(r.Context(), w)
	}
}
