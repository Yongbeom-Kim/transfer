package middleware

import (
	"net/http"
	"regexp"
)

var ALLOWED_ORIGINS = []*regexp.Regexp{
	regexp.MustCompile("http://localhost:.*"),
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			for _, origin := range ALLOWED_ORIGINS {
				if origin.MatchString(r.Header.Get("Origin")) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
					w.WriteHeader(http.StatusOK)
					return
				}
			}
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
