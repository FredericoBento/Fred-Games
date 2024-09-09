package handler

import "net/http"

func IsHTMX(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	} else {
		return false
	}
}
