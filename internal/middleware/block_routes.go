package middleware

import (
	"net/http"
	"strings"
)

var (
	blockedRoutesPrefixes = make([]string, 0)
)

func BlockRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, prefix := range blockedRoutesPrefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				http.Error(w, "Service temporarily unavailable :(", http.StatusServiceUnavailable)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func BlockRoute(prefix string) {
	blockedRoutesPrefixes = append(blockedRoutesPrefixes, prefix)
}

func UnblockRoute(prefix string) {
	for i, route := range blockedRoutesPrefixes {
		if route == prefix {
			blockedRoutesPrefixes = append(blockedRoutesPrefixes[:i], blockedRoutesPrefixes[i+1:]...)
		}
	}
}
