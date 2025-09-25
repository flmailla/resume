package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/flmailla/resume/models"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		sentHeader       string
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "No Authorization Header",
			sentHeader: "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: models.ErrNoTokenSent.Error(),
		},
		{
			name: "No Bearer in Authorization Header",
			sentHeader: "xxxx",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: models.ErrNotBearer.Error(),
		},
		{
			name: "Token is empty",
			sentHeader: "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: models.ErrUnauthorized.Error(),
		},
		{
			name: "With Authorization Header but invalid token",
			sentHeader: "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &JWTValidator{}
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})
			middleware := validator.AuthMiddleware(nextHandler)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.sentHeader != "" {
				req.Header.Set("Authorization", tt.sentHeader)
			}
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if rr.Code != http.StatusUnauthorized {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
