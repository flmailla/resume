package auth

import (
	"net/http"
	"strings"
	"github.com/flmailla/resume/models"
)

// Middleware used by net/http
// used to check the request authorization
// and redirect to the right MUX handler afterwards
func (v *JWTValidator) AuthMiddleware(mux http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			mux.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)

		if authHeader == "" {
			w.Write([]byte(models.ErrNoTokenSent.Error()))
            return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			w.Write([]byte(models.ErrNotBearer.Error()))
            return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		if token == "" {
			w.Write([]byte(models.ErrUnauthorized.Error()))
            return
		}

		err := v.verifyToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		mux.ServeHTTP(w, r)
	})
}