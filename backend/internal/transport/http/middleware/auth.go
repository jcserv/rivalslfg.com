package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jcserv/rivalslfg/internal/auth"
)

// Change middleware to work with http.HandlerFunc instead of http.Handler
func RequireRight(right auth.Right) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value("claims").(jwt.MapClaims)

			if rights, ok := claims["rights"].([]string); ok {
				hasRight := false
				for _, r := range rights {
					if auth.IsEqual(r, right) {
						hasRight = true
						break
					}
				}

				if !hasRight {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			next(w, r)
		}
	}
}
