package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		sentHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No Authorization Header",
			sentHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   ascii401,
		},
		{
			name:           "No Bearer in Authorization Header",
			sentHeader:     "xxxx",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   ascii401,
		},
		{
			name:           "Token is empty",
			sentHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   ascii401,
		},
		{
			name:           "With Authorization Header but invalid token",
			sentHeader:     "Bearer invalid",
			expectedStatus: http.StatusForbidden,
			expectedBody:   ascii403,
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

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestAuthMiddlewareSpecificRoute(t *testing.T) {
	tests := []struct {
		name           string
		sentHeader     string
		expectedStatus int
	}{
		{
			name:           "No token needed",
			sentHeader:     "",
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &JWTValidator{}
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})
			middleware := validator.AuthMiddleware(nextHandler)

			req := httptest.NewRequest("GET", "/health", nil)
			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

		})
	}
}
