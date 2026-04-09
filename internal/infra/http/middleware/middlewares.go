package middleware

import "net/http"

type TokenService interface {
	Validate(tokenString string) (string, error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
