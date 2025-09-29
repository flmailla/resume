package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealth(t *testing.T) {
	tests := []struct {
		name              string
		mockStore         *mockStore
		wantStatusCode    int
		wantStatusMessage string
	}{
		{
			name:              "Retrieve status",
			mockStore:         nil,
			wantStatusCode:    http.StatusOK,
			wantStatusMessage: "healthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			healthHandler := NewHealthHandler(tt.mockStore)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /health", healthHandler.GetHealthStatus)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/health", nil)

			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}

			var got responseStatus

			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("%v", w.Body)
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if got.Status != tt.wantStatusMessage {
				t.Errorf("expected error message %q, got %q", tt.wantStatusMessage, got.Status)
			}
		})
	}
}
