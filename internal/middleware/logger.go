package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{w, http.StatusOK}
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := NewStatusRecorder(w)

		next.ServeHTTP(rec, r)

		duration := time.Since(start).String()

		logString := duration + " " + r.Method + " " + r.URL.Path + " - " + strconv.Itoa(rec.statusCode)
		slog.Info(logString)

		// If there are query parameters, log them
		queryParams := r.URL.Query()
		if len(queryParams) > 0 {
			// Start logging query parameters
			queryLogString := " - Query Parameters:"
			for key, values := range queryParams {
				queryLogString += "\n\t\t\t\t\t - " + key + ":"
				for _, val := range values {
					queryLogString += " " + val
				}
			}
			slog.Info(queryLogString)
		}
	})
}
