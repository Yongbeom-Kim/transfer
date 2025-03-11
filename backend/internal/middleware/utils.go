package middleware

import "net/http"

// Compose a list of middlewares into a single middleware
// First middleware is executed first, then second, etc.
func Compose(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for _, middleware := range middlewares {
			next = middleware(next)
		}
		return next
	}
}
