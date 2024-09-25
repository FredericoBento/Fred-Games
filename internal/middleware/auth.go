package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/services"
)

var (
	authService *services.AuthService
)

func SetAuthService(service *services.AuthService) {
	authService = service
}

func RequiredLogged(next http.Handler) http.Handler {
	if authService == nil {
		log.Fatal("authservice not setup")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie(authService.GetCookieName())
		if err != nil {
			// http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		_, ok := authService.ValidateSession(context.TODO(), c.Value)
		if ok != nil {
			// http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequiredAdmin(next http.Handler) http.Handler {
	if authService == nil {
		log.Fatal("authservice not setup")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie(authService.GetCookieName())
		if err != nil {
			// http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		user, ok := authService.ValidateSession(context.TODO(), c.Value)
		if ok != nil {
			// http.Error(w, "Forbidden, invalid token", http.StatusForbidden)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		if user.Username == "fred" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden, not an admin", http.StatusForbidden)
			return
		}

	})
}
