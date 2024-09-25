package handler

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/auth_views"
	"github.com/a-h/templ"
)

type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
	log         *slog.Logger
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	if authService == nil {
		log.Fatal("auth service not provided")
	}
	if userService == nil {
		log.Fatal("user service not provided")
	}
	lo, err := logger.NewHandlerLogger("AuthHandler", "", true)
	if err != nil {
		lo = slog.Default()
	}

	return &AuthHandler{
		authService: authService,
		userService: userService,
		log:         lo,
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/sign-in":
		ah.signIn(w, r)
	case "/sign-up":
		ah.signUp(w, r)
	case "/logout":
		ah.GetLogout(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	}
}

func (ah *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ah.GetSignIn(w, r)
	case http.MethodPost:
		ah.PostSignIn(w, r)
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

	token, err := ah.authService.GetToken(r)
	if err == nil && token != "" {
		user, err := ah.authService.ValidateSession(context.TODO(), token)
		if err == nil && user != nil {
			if ah.authService.IsAdmin(user.Username) {
				// http.Redirect(w, r, "/admin/", http.StatusSeeOther)
				Redirect(w, r, "/admin/")
				return
			} else {
				// http.Redirect(w, r, "/home/", http.StatusSeeOther)
				Redirect(w, r, "/home/")
				return
			}
		}
	}

	ah.View(w, r, viewProps{
		title:   "Sign Up",
		content: auth_views.SignUpForm(),
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

	// w.WriteHeader(http.StatusCreated)
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)

}

func (ah *AuthHandler) GetSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if ah.authService == nil {
		ah.log.Error("NO AUTH SERVICE")
		return
	}

	token, err := ah.authService.GetToken(r)
	if err == nil && token != "" {
		user, err := ah.authService.ValidateSession(context.TODO(), token)
		if err == nil && user != nil {
			if ah.authService.IsAdmin(user.Username) {
				// http.Redirect(w, r, "/admin/", http.StatusSeeOther)
				Redirect(w, r, "/admin/")
				return
			} else {
				// http.Redirect(w, r, "/home/", http.StatusSeeOther)
				Redirect(w, r, "/home/")
				return
			}
		}
	}

	ah.View(w, r, viewProps{
		title:   "Sign In",
		content: auth_views.SignInForm(),
	})
}

func (ah *AuthHandler) GetLogout(w http.ResponseWriter, r *http.Request) {
	token, err := ah.authService.GetToken(r)
	if err != nil {
		http.Error(w, "No token, could not logout", http.StatusBadRequest)
		return
	}

	c, err := r.Cookie(ah.authService.GetCookieName())
	if err != http.ErrNoCookie && c != nil {
		c.Value = ""
		c.Expires = time.Now()
		http.SetCookie(w, c)
	} else {
		// delete with other method in case of no cookie auth
	}

	ah.authService.DestroySession(context.Background(), token)
	Redirect(w, r, "/sign-in")
}

func (ah *AuthHandler) PostSignIn(w http.ResponseWriter, r *http.Request) {
	type SignInForm struct {
		username string
		password string
	}

	data := SignInForm{}
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")

	if data.username == "" || data.password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := ah.userService.UserExists(context.Background(), data.username)
	if err != nil {
		ah.log.Error(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	if !exists {
		ah.log.Error("user does not exist")
		w.Write([]byte("user does not exist"))
		return
	}

	u, err := ah.authService.Authenticate(context.TODO(), data.username, data.password)
	if err != nil {
		ah.log.Error(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	if u != nil {
		token, err := ah.authService.CreateSession(u)
		if err != nil {
			ah.log.Error(err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		cookie := http.Cookie{
			Name:  "session_token",
			Value: token,
		}

		http.SetCookie(w, &cookie)

		if ah.authService.IsAdmin(u.Username) {
			Redirect(w, r, "/admin/")
			return
		}
		Redirect(w, r, "/home/")
		return
	}
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
