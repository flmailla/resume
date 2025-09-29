package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test the writeJSON function
func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		payload        interface{}
		expectedStatus int
		expectedBody   string
		expectError    bool
	}{
		{
			name:           "successful JSON response",
			status:         http.StatusOK,
			payload:        map[string]string{"message": "hello"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"hello"}`,
			expectError:    false,
		},
		{
			name:           "empty payload",
			status:         http.StatusCreated,
			payload:        nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   "null",
			expectError:    false,
		},
		{
			name:           "struct payload",
			status:         http.StatusAccepted,
			payload:        struct{ ID int }{ID: 123},
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"ID":123}`,
			expectError:    false,
		},
		{
			name:           "array payload",
			status:         http.StatusOK,
			payload:        []string{"a", "b", "c"},
			expectedStatus: http.StatusOK,
			expectedBody:   `["a","b","c"]`,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the function under test
			writeJSON(rr, tt.status, tt.payload)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("writeJSON() status = %v, want %v", rr.Code, tt.expectedStatus)
			}

			// Check Content-Type header
			expectedContentType := "application/json"
			if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("writeJSON() Content-Type = %v, want %v", contentType, expectedContentType)
			}

			// Check response body (trim newline added by json.Encoder)
			body := strings.TrimSpace(rr.Body.String())
			if body != tt.expectedBody {
				t.Errorf("writeJSON() body = %v, want %v", body, tt.expectedBody)
			}
		})
	}
}

func TestWriteJSONInvalidPayload(t *testing.T) {
	rr := httptest.NewRecorder()

	// Channels cannot be marshaled to JSON
	invalidPayload := make(chan int)

	writeJSON(rr, http.StatusOK, invalidPayload)

	// The function should handle the error and return 500
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("writeJSON() with invalid payload status = %v, want %v", rr.Code, http.StatusInternalServerError)
	}

	body := strings.TrimSpace(rr.Body.String())
	if body != "Internal Server Error" {
		t.Errorf("writeJSON() with invalid payload body = %v, want %v", body, "Internal Server Error")
	}
}
