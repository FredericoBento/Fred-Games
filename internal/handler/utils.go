package handler

import (
	"net/http"

	"github.com/FredericoBento/HandGame/internal/middleware"
	"github.com/FredericoBento/HandGame/internal/models"
)

func IsHTMX(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	} else {
		return false
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, route string) {
	if IsHTMX(r) {
		w.Header().Add("HX-Redirect", route)
	} else {
		http.Redirect(w, r, route, http.StatusSeeOther)
	}
}

func GetLoggedUser(r *http.Request) (*models.User, bool) {
	user := r.Context().Value(middleware.LoggedUserKey)
	if user != nil {
		return user.(*models.User), true
	}
	return nil, false
}

func IsLogged(r *http.Request) bool {
	_, logged := GetLoggedUser(r)
	return logged
}

func IsAdmin(r *http.Request) bool {
	isAdmin := r.Context().Value(middleware.IsAdminKey)
	if isAdmin != nil {
		return isAdmin.(bool)
	}
	return false
}
