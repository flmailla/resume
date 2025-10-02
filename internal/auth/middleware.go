package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware used by net/http
// used to check the request authorization
// and redirect to the right MUX handler afterwards
func (v *JWTValidator) AuthMiddleware(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			mux.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, ascii401)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, ascii401)
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		if token == "" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, ascii401)
			return
		}

		err := v.verifyToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, ascii403)
			return
		}

		mux.ServeHTTP(w, r)
	})
}
