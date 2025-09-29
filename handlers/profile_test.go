package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/flmailla/resume/models"
)

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             *models.Profile
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: "unknown error",
			mockStore: &mockStore{
				GetProfileFunc: func(profileId int) (*models.Profile, error) {
					return &models.Profile{}, errors.New("unknown error")
				},
			},
			want:             &models.Profile{},
			wantStatusCode:   http.StatusInternalServerError,
			wantErrorMessage: models.ErrProfileNotFetched.Error(),
		},
		{
			name: "successful query with multiple profiles",
			mockStore: &mockStore{
				GetProfileFunc: func(profileId int) (*models.Profile, error) {
					return &models.Profile{
							ID:         1,
							FirstName:  "FN1",
							LastName:   "LN1",
							BirthDate:  time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
							Pronoun:    "He",
							Email:      "email@maillard.ch",
							Location:   "Switzerland",
							PostalCode: 1000,
							Headline:   "Headline",
							About:      "about",
						},
						nil
				},
			},
			want: &models.Profile{
				ID:         1,
				FirstName:  "FN1",
				LastName:   "LN1",
				BirthDate:  time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
				Pronoun:    "He",
				Email:      "email@maillard.ch",
				Location:   "Switzerland",
				PostalCode: 1000,
				Headline:   "Headline",
				About:      "about",
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileHandler := NewProfileHandler(tt.mockStore)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /profiles/{profile_id}", profileHandler.GetProfile)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/profiles/1", nil)

			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", w.Code, tt.wantStatusCode)
			}

			if w.Code == http.StatusOK {
				var got *models.Profile

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Store.GetProfile() got user %+v, want %+v", got, tt.want)
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

func TestGetProfileWrongPathParameter(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             *models.Profile
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrInvalidId.Error(),
			mockStore: &mockStore{
				GetProfileFunc: func(profileId int) (*models.Profile, error) {
					return &models.Profile{}, models.ErrInvalidId
				},
			},
			want:             &models.Profile{},
			wantStatusCode:   http.StatusInternalServerError,
			wantErrorMessage: models.ErrInvalidId.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileHandler := NewProfileHandler(tt.mockStore)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /profiles/{profile_id}", profileHandler.GetProfile)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/profiles/abc", nil)

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
