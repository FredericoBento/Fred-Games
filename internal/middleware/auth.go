package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/FredericoBento/HandGame/internal/models"
)

var (
	authService AuthService
)

type contextKey string

const (
	LoggedUserKey = contextKey("user")
)

type AuthService interface {
	GetCookieName() string
	ValidateSession(ctx context.Context, token string) (*models.User, error)
}

func SetAuthService(service AuthService) {
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
	next = AddUserToContext(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := GetUserFromContext(r)
		if !ok {
			// http.Error(w, "Forbidden", http.StatusForbidden)
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

func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user := r.Context().Value(LoggedUserKey)
	if user != nil {
		return user.(*models.User), true
	}
	return nil, false
}

// Will add the user to context if such is logged
func AddUserToContext(next http.Handler) http.Handler {
	if authService == nil {
		log.Fatal("authService needs to be provided to use this middleware")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(authService.GetCookieName())
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, ok := authService.ValidateSession(context.TODO(), c.Value)
		if ok != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), LoggedUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
