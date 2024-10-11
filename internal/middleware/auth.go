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
	IsAdminKey    = contextKey("isAdmin")
)

type AuthService interface {
	GetToken(r *http.Request) (string, error)
	GetCookieName() string
	ValidateSession(ctx context.Context, token string) (*models.User, error)
	IsAdmin(username string) bool
}

func SetAuthService(service AuthService) {
	authService = service
}

func RequiredLogged(next http.Handler) http.Handler {
	if authService == nil {
		log.Fatal("authservice not setup")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(LoggedUserKey)
		if user != nil {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return

	})
}

func RequiredAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(LoggedUserKey)
		if user == nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		isAdmin := r.Context().Value(IsAdminKey)
		if isAdmin != nil {
			if isAdmin.(bool) {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Forbidden, not an admin", http.StatusForbidden)
		return
	})
}

// Will Add Logged User and Admin level to context
func AuthEssential(next http.Handler) http.Handler {
	if authService == nil {
		log.Fatal("authService needs to be provided to use this middleware")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(addUserToContext(r))
		r = r.WithContext(addAdminLevelToContext(r.Context()))
		next.ServeHTTP(w, r)
		return
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
		return
	})

}

func addAdminLevelToContext(ctx context.Context) context.Context {
	if authService == nil {
		log.Fatal("authService needs to be provided to use this middleware")
	}
	user := ctx.Value(LoggedUserKey)
	if user != nil {
		if authService.IsAdmin(user.(*models.User).Username) {
			return context.WithValue(ctx, IsAdminKey, true)
		}
	}
	return context.WithValue(ctx, IsAdminKey, false)
}

func addUserToContext(r *http.Request) context.Context {
	if authService == nil {
		log.Fatal("authService needs to be provided to use this middleware")
	}
	token, err := authService.GetToken(r)
	if err != nil {
		return r.Context()
	}

	user, err := authService.ValidateSession(r.Context(), token)
	if err != nil {
		return context.WithValue(r.Context(), LoggedUserKey, nil)
	}

	return context.WithValue(r.Context(), LoggedUserKey, user)
}
