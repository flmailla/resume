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

func TestGetExperiencesByProfile(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Experience
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: "unknown error",
			mockStore: &mockStore{
				GetDistinctExperiencesByProfileFunc: func(profileId int) ([]models.Experience, error)	 {
					return []models.Experience{}, errors.New("unknown error")
				},
			},
			want: []models.Experience{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrExperiencesNotFetched.Error(),
		},
		{
			name: "successful query with multiple educations",
			mockStore: &mockStore{
				GetDistinctExperiencesByProfileFunc: func(profileId int) ([]models.Experience, error)	 {
					return []models.Experience{
							{
								ID: 1, 
								Title: "Job1", 
								Company: "Company1", 
								StartDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
								EndDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
								Location: "Switzerland",
								Description: "A super job",
							},
							{
								ID: 2, 
								Title: "Job2", 
								Company: "Company3", 
								StartDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
								EndDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
								Location: "Switzerland",
								Description: "A terrific one",
							},
						},
					nil
				},
			},
			want: []models.Experience{
				{
					ID: 1,
					Title: "Job1", 
					Company: "Company1", 
					StartDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					EndDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					Location: "Switzerland",
					Description: "A super job",
				},
				{
					ID: 2, 
					Title: "Job2", 
					Company: "Company3", 
					StartDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					EndDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					Location: "Switzerland",
					Description: "A terrific one",
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			educationHandler := NewExperienceHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/experiences", educationHandler.GetExperiencesByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/1/experiences", nil)
			
			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", w.Code, tt.wantStatusCode)
			}

			if w.Code == http.StatusOK {
				var got []models.Experience

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				for i, experience := range got {
					if i < len(tt.want) {
						expected := tt.want[i]
						if experience.ID != expected.ID {
							t.Errorf("Experience[%d].ID = %v, want %v", i, experience.ID, expected.ID)
						}
						if experience.Title != expected.Title {
							t.Errorf("Experience[%d].Title = %v, want %v", i, experience.Title, expected.Title)
						}
						if experience.Company != expected.Company {
							t.Errorf("Experience[%d].Company = %v, want %v", i, experience.Company, expected.Company)
						}
						if !experience.StartDate.Equal(expected.StartDate) {
							t.Errorf("Experience[%d].StartDate = %v, want %v", i, experience.StartDate, expected.StartDate)
						}
						if !experience.EndDate.Equal(expected.EndDate) {
							t.Errorf("Experience[%d].EndDate = %v, want %v", i, experience.EndDate, expected.EndDate)
						}
						if experience.Location != expected.Location {
							t.Errorf("Experience[%d].Location = %v, want %v", i, experience.Location, expected.Location)
						}
						if experience.Description != expected.Description {
							t.Errorf("Experience[%d].Description = %v, want %v", i, experience.Description, expected.Description)
						}
						if !models.SkillsEqual(experience.Skills, expected.Skills) {
							t.Errorf("Experience[%d].Skills = %v, want %v", i, experience.Skills, expected.Skills)
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

func TestGetExperiencesByProfileWrongPathParameter(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Experience
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrInvalidId.Error(),
			mockStore: &mockStore{
				GetDistinctExperiencesByProfileFunc: func(profileId int) ([]models.Experience, error)	 {
					return []models.Experience{}, models.ErrInvalidId
				},
			},
			want: []models.Experience{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrInvalidId.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			experienceHandler := NewExperienceHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/experiences", experienceHandler.GetExperiencesByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/abc/experiences", nil)
			
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