package handlers

import (
	"time"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/flmailla/resume/models"
)

func TestGetEducationsByProfile(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Education
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: "unknown error",
			mockStore: &mockStore{
				GetDistinctEducationsByProfileFunc: func(profileId int) ([]models.Education, error)	 {
					return []models.Education{}, errors.New("unknown error")
				},
			},
			want: []models.Education{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrEducationsNotFetched.Error(),
		},
		{
			name: "successful query with multiple educations",
			mockStore: &mockStore{
				GetDistinctEducationsByProfileFunc: func(profileId int) ([]models.Education, error)	 {
					return []models.Education{
							{ID: 1, Title: "University", Issued: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), Description: "A university journey"},
							{ID: 2, Title: "University again", Issued: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), Description: "Another university journey"},
						},
					nil
				},
			},
			want: []models.Education{
				{ID: 1, Title: "University", Issued: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), Description: "A university journey"},
				{ID: 2, Title: "University again", Issued: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), Description: "Another university journey"},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			educationHandler := NewEducationHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/educations", educationHandler.GetEducationsByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/1/educations", nil)
			
			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", w.Code, tt.wantStatusCode)
			}

			if w.Code == http.StatusOK {
				var got []models.Education

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				for i, education := range got {
					if i < len(tt.want) {
						if education != tt.want[i] {
							t.Errorf("Store.GetDistinctEducationsByProfile() got user %+v, want %+v", education, tt.want[i])
						}
					}
				}

			} else {
				var got map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				if got["error"] != tt.wantErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.wantErrorMessage, got["error"])
				}
			}

		})
	}
}


func TestGetEducationsByProfileWrongPathParameter(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Education
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrInvalidId.Error(),
			mockStore: &mockStore{
				GetDistinctEducationsByProfileFunc: func(profileId int) ([]models.Education, error)	 {
					return []models.Education{}, models.ErrInvalidId
				},
			},
			want: []models.Education{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrInvalidId.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			educationHandler := NewEducationHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/educations", educationHandler.GetEducationsByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/abc/educations", nil)
			
			mux.ServeHTTP(w, r)

			var got map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("%v", w.Body)
				t.Fatalf(models.ErrUnmarshal.Error(), err)
			}

			if got["error"] != tt.wantErrorMessage {
				t.Errorf("expected error message %q, got %q", tt.wantErrorMessage, got["error"])
			}
		
		})
	}
}