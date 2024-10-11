package handler

import (
	"context"
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
	lo, err := logger.NewHandlerLogger("AuthHandler", "", false)
	if err != nil {
		lo = slog.Default()
		lo.Error("could not create admin handler logger, using slog.default")
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

	// token, err := ah.authService.GetToken(r)
	// if err == nil && token != "" {
	// user, err := ah.authService.ValidateSession(context.TODO(), token)
	// if err == nil && user != nil {
	if IsLogged(r) {
		if IsAdmin(r) {
			// if ah.authService.IsAdmin(user.Username) {
			// http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			Redirect(w, r, "/admin")
			return
		}
		// } else {
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
		Redirect(w, r, "/home")
		return
	}

	ah.returnSignUpForm(w, r, auth_views.SignUpFormData{})
	return
}

func (ah *AuthHandler) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := auth_views.SignUpFormData{}
	data.Username = r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("repeat_password")

	errFound := false

	if data.Username == "" {
		data.UsernameErr = "This field cannot be empty"
		errFound = true
	}

	if password == "" {
		data.PasswordErr = "This field cannot be empty"
		errFound = true
	}

	if confirmPassword == "" {
		data.ConfirmPasswordErr = "This field cannot be empty"
		errFound = true
	}

	if confirmPassword != password {
		data.PasswordErr = "Passwords do not match"
		data.ConfirmPasswordErr = "Passwords do not match"
		errFound = true
	}

	if errFound {
		ah.returnSignUpForm(w, r, data)
		return
	}

	err := ah.userService.CreateUser(context.Background(), &models.User{ID: 0, Username: data.Username, Password: password})
	if err != nil {
		ah.log.Error(err.Error())
		switch err {
		case services.ErrCouldNotCreateUser:
			w.WriteHeader(http.StatusInternalServerError)
			data.GeneralErr = "A server error has ocurred, try again later"

		case services.ErrUserAlreadyExists:
			w.WriteHeader(http.StatusBadRequest)
			data.UsernameErr = "This username is already taken"

		case services.ErrUserExistsFailed:
			w.WriteHeader(http.StatusInternalServerError)
			data.GeneralErr = "A server error has ocurred, try again later"
		}

		ah.returnSignUpForm(w, r, data)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	Redirect(w, r, "/sign-in")
}

func (ah *AuthHandler) GetSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if ah.authService == nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	data := auth_views.SignInFormData{
	// 		GeneralErr: "A server error has ocurred, try again later",
	// 	}
	// 	ah.returnSignInForm(w, r, data)
	// 	return
	// }

	if IsLogged(r) {
		if IsAdmin(r) {
			Redirect(w, r, "/admin")
			return
		}
		Redirect(w, r, "/home")
		return
	}

	ah.returnSignInForm(w, r, auth_views.SignInFormData{})
	return
}

func (ah *AuthHandler) PostSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data := auth_views.SignInFormData{}

	data.Username = r.FormValue("username")
	password := r.FormValue("password")

	data.UsernameErr = ""
	data.PasswordErr = ""
	data.GeneralErr = ""

	if data.Username == "" {
		data.UsernameErr = "This username is invalid"
	}

	if password == "" {
		data.PasswordErr = "This password is invalid"
	}

	if data.UsernameErr != "" || data.PasswordErr != "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.returnSignInForm(w, r, data)
		return
	}

	exists, err := ah.userService.UserExists(context.Background(), data.Username)
	if err != nil {
		ah.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		data.GeneralErr = "A server error ocurred, try again later"
		ah.returnSignInForm(w, r, data)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		data.UsernameErr = "This user does not exist"
		ah.returnSignInForm(w, r, data)
		return
	}

	u, err := ah.authService.Authenticate(context.TODO(), data.Username, password)
	if err != nil {
		ah.log.Error(err.Error())
		switch err {
		case services.ErrIncorrectCredentials:
			w.WriteHeader(http.StatusBadRequest)
			data.GeneralErr = "Incorrect Credentials"

		case services.ErrCouldNotFindUser:
			w.WriteHeader(http.StatusBadRequest)
			data.UsernameErr = "This user was not found"

		default:
			w.WriteHeader(http.StatusInternalServerError)
			data.GeneralErr = "A server error ocurred, try again later"
		}

		ah.returnSignInForm(w, r, data)
		return
	}

	if u != nil {
		token, err := ah.authService.CreateSession(u)
		if err != nil {
			ah.log.Error(err.Error())
			data.GeneralErr = "A server error ocurred, try again later"
			ah.returnSignInForm(w, r, data)
			return
		}

		cookie := http.Cookie{
			Name:  "session_token",
			Value: token,
		}

		http.SetCookie(w, &cookie)

		if ah.authService.IsAdmin(u.Username) {
			Redirect(w, r, "/admin")
			return
		}
		Redirect(w, r, "/home")
		return
	}
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

type ViewAuthProps struct {
	title   string
	content templ.Component
}

func (ah *AuthHandler) View(w http.ResponseWriter, r *http.Request, props ViewAuthProps) {
	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, auth_views.DefaultNavbar(), props.content).Render(r.Context(), w)
	}
}

func (ah *AuthHandler) returnSignInForm(w http.ResponseWriter, r *http.Request, data auth_views.SignInFormData) {

	ah.View(w, r, ViewAuthProps{
		content: auth_views.SignInForm(data),
	})
}

func (ah *AuthHandler) returnSignUpForm(w http.ResponseWriter, r *http.Request, data auth_views.SignUpFormData) {
	ah.View(w, r, ViewAuthProps{
		content: auth_views.SignUpForm(data),
	})

}
