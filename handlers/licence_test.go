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

func TestGetLicencesByProfile(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Licence
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: "unknown error",
			mockStore: &mockStore{
				GetDistinctLicencesByProfileFunc: func(profileId int) ([]models.Licence, error)	 {
					return []models.Licence{}, errors.New("unknown error")
				},
			},
			want: []models.Licence{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrLicencesNotFetched.Error(),
		},
		{
			name: "successful query with multiple licences",
			mockStore: &mockStore{
				GetDistinctLicencesByProfileFunc: func(profileId int) ([]models.Licence, error)	 {
					return []models.Licence{
							{
								ID: 1, 
								Title: "Licence1", 
								Issuer: "Issuer1", 
								IssuedAt: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
								LicenceType: models.CERTIFICATION,
							},
							{
								ID: 2, 
								Title: "Licence2", 
								Issuer: "Issuer2", 
								IssuedAt: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
								LicenceType: models.LICENCE,
							},
						},
					nil
				},
			},
			want: []models.Licence{
				{
					ID: 1, 
					Title: "Licence1", 
					Issuer: "Issuer1", 
					IssuedAt: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					LicenceType: models.CERTIFICATION,
				},
				{
					ID: 2, 
					Title: "Licence2", 
					Issuer: "Issuer2", 
					IssuedAt: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
					LicenceType: models.LICENCE,
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			licenceHandler := NewLicenceHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/licences", licenceHandler.GetLicencesByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/1/licences", nil)
			
			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", w.Code, tt.wantStatusCode)
			}

			if w.Code == http.StatusOK {
				var got []models.Licence

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				for i, licence := range got {
					if i < len(tt.want) {
						if licence != tt.want[i] {
							t.Errorf("Store.GetDistinctLicencesByProfile() got user %+v, want %+v", licence, tt.want[i])
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

func TestGetLicencesByProfileWrongPathParameter(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Licence
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrInvalidId.Error(),
			mockStore: &mockStore{
				GetDistinctLicencesByProfileFunc: func(profileId int) ([]models.Licence, error)	 {
					return []models.Licence{}, models.ErrInvalidId
				},
			},
			want: []models.Licence{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrInvalidId.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			licenceHandler := NewLicenceHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/licences", licenceHandler.GetLicencesByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/abc/licences", nil)
			
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