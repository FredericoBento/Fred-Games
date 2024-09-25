package handler

import "net/http"

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
