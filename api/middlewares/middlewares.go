package middlewares

import (
	"errors"
	"net/http"
	"github.com/haikalvidya/goApiBlog/api/responses"
	"github.com/haikalvidya/goApiBlog/api/auth"
)

// middleware for json responses
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w,r)
	}
}

// middleware for authetication
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		}
		next(w,r)
	}
}