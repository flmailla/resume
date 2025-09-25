package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/flmailla/resume/models"
)

func TestGetSkills(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Skill
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrUnknown.Error(),
			mockStore: &mockStore{
				GetDistinctSkillsFunc: func() ([]models.Skill, error)	 {
					return []models.Skill{}, models.ErrUnknown
				},
			},
			want: []models.Skill{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrSkillsNotFetched.Error(),
		},
		{
			name: "successful query with multiple skills",
			mockStore: &mockStore{
				GetDistinctSkillsFunc: func() ([]models.Skill, error)	 {
					return []models.Skill{
							{
								ID: 1, 
								Name: "Skill1",
							},
							{
								ID: 2, 
								Name: "Skill2",
							},
						},
					nil
				},
			},
			want: []models.Skill{
				{
					ID: 1, 
					Name: "Skill1",
				},
				{
					ID: 2, 
					Name: "Skill2",
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillHandler := NewSkillHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /skills", skillHandler.GetSkills)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/skills", nil)
			
			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}

			if w.Code == http.StatusOK {
				var got []models.Skill

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf(models.ErrUnmarshal.Error(), err)
				}

				for i, skill := range got {
					if i < len(tt.want) {
						if skill != tt.want[i] {
							t.Errorf("Store.GetSkills() got skill %+v, want %+v", skill, tt.want[i])
						}
					}
				}

			} else {
				var got map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf(models.ErrUnmarshal.Error(), err)
				}

				if got["error"] != tt.wantErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.wantErrorMessage, got["error"])
				}
			}

		})
	}
}

func TestGetSkillsByProfile(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Skill
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrUnknown.Error(),
			mockStore: &mockStore{
				GetDistinctSkillsByProfileFunc: func(profileId int) ([]models.Skill, error)	 {
					return []models.Skill{}, models.ErrUnknown
				},
			},
			want: []models.Skill{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrSkillsNotFetched.Error(),
		},
		{
			name: "successful query with multiple skills by Profile",
			mockStore: &mockStore{
				GetDistinctSkillsByProfileFunc: func(profileId int) ([]models.Skill, error)	 {
					return []models.Skill{
							{
								ID: 1, 
								Name: "Skill1",
							},
							{
								ID: 2, 
								Name: "Skill2",
							},
						},
					nil
				},
			},
			want: []models.Skill{
				{
					ID: 1, 
					Name: "Skill1",
				},
				{
					ID: 2, 
					Name: "Skill2",
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillHandler := NewSkillHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/skills", skillHandler.GetSkillsByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/1/skills", nil)
			
			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}

			if w.Code == http.StatusOK {
				var got []models.Skill

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf(models.ErrUnmarshal.Error(), err)
				}

				for i, licence := range got {
					if i < len(tt.want) {
						if licence != tt.want[i] {
							t.Errorf("Store.GetSkillsByProfile() got skill %+v, want %+v", licence, tt.want[i])
						}
					}
				}

			} else {
				var got map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf(models.ErrUnmarshal.Error(), err)
				}

				if got["error"] != tt.wantErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.wantErrorMessage, got["error"])
				}
			}

		})
	}
}

func TestGetSkillsByExperience(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Skill
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrUnknown.Error(),
			mockStore: &mockStore{
				GetDistinctSkillsByExperienceFunc: func(experienceId int) ([]models.Skill, error)	 {
					return []models.Skill{}, models.ErrUnknown
				},
			},
			want: []models.Skill{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrSkillsNotFetched.Error(),
		},
		{
			name: "successful query with multiple skills by experience",
			mockStore: &mockStore{
				GetDistinctSkillsByExperienceFunc: func(experienceId int) ([]models.Skill, error)	 {
					return []models.Skill{
							{
								ID: 1, 
								Name: "Skill1",
							},
							{
								ID: 2, 
								Name: "Skill2",
							},
						},
					nil
				},
			},
			want: []models.Skill{
				{
					ID: 1, 
					Name: "Skill1",
				},
				{
					ID: 2, 
					Name: "Skill2",
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillHandler := NewSkillHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{experience_id}/skills", skillHandler.GetSkillsByExperience)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/1/skills", nil)
			
			mux.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d",tt.wantStatusCode, w.Code)
			}

			if w.Code == http.StatusOK {
				var got []models.Skill

				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf(models.ErrUnmarshal.Error(), err)
				}

				for i, licence := range got {
					if i < len(tt.want) {
						if licence != tt.want[i] {
							t.Errorf("Store.GetSkillsByExperience() got skill %+v, want %+v", licence, tt.want[i])
						}
					}
				}

			} else {
				var got map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("%v", w.Body)
					t.Fatalf(models.ErrUnmarshal.Error(), err)
				}

				if got["error"] != tt.wantErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.wantErrorMessage, got["error"])
				}
			}

		})
	}
}

func TestGetSkillsByExperienceWrongPathParameter(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Skill
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrInvalidId.Error(),
			mockStore: &mockStore{
				GetDistinctSkillsByExperienceFunc: func(experienceId int) ([]models.Skill, error)	 {
					return []models.Skill{}, models.ErrInvalidId
				},
			},
			want: []models.Skill{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrInvalidId.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillHandler := NewSkillHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{experience_id}/skills", skillHandler.GetSkillsByExperience)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/abc/skills", nil)
			
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

func TestGetSkillsByProfileWrongPathParameter(t *testing.T) {
	tests := []struct {
		name             string
		mockStore        *mockStore
		want             []models.Skill
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name: models.ErrInvalidId.Error(),
			mockStore: &mockStore{
				GetDistinctSkillsByProfileFunc: func(profileId int) ([]models.Skill, error)	 {
					return []models.Skill{}, models.ErrInvalidId
				},
			},
			want: []models.Skill{},
			wantStatusCode: http.StatusInternalServerError,
			wantErrorMessage: models.ErrInvalidId.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillHandler := NewSkillHandler(tt.mockStore)

			mux := http.NewServeMux()
    		mux.HandleFunc("GET /profiles/{profile_id}/skills", skillHandler.GetSkillsByProfile)

			w := httptest.NewRecorder()
		    r := httptest.NewRequest("GET", "/profiles/abc/skills", nil)
			
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

func TestSkillsEqual(t *testing.T) {
	tests := []struct {
		name     string
		skillsA  []models.Skill
		skillsB  []models.Skill
		expected bool
	}{
		{
			name:     "both nil slices",
			skillsA:  nil,
			skillsB:  nil,
			expected: true,
		},
		{
			name:     "both empty initialized slices",
			skillsA:  make([]models.Skill, 0),
			skillsB:  make([]models.Skill, 0),
			expected: true,
		},
		{
			name:     "one nil, one empty",
			skillsA:  nil,
			skillsB:  make([]models.Skill, 0),
			expected: true,
		},
		{
			name:     "identical single elements",
			skillsA:  []models.Skill{{ID: 1, Name: "git"}},
			skillsB:  []models.Skill{{ID: 1, Name: "git"}},
			expected: true,
		},
		{
			name:     "different IDs",
			skillsA:  []models.Skill{{ID: 1, Name: "git"}},
			skillsB:  []models.Skill{{ID: 2, Name: "git"}},
			expected: false,
		},
		{
			name:     "different names",
			skillsA:  []models.Skill{{ID: 1, Name: "git"}},
			skillsB:  []models.Skill{{ID: 1, Name: "docker"}},
			expected: false,
		},
		{
			name:     "different lengths",
			skillsA:  []models.Skill{{ID: 1, Name: "git"}},
			skillsB:  []models.Skill{{ID: 1, Name: "git"}, {ID: 2, Name: "docker"}},
			expected: false,
		},
		{
			name: "multiple identical elements",
			skillsA: []models.Skill{
				{ID: 1, Name: "git"},
				{ID: 2, Name: "docker"},
				{ID: 3, Name: "kubernetes"},
			},
			skillsB: []models.Skill{
				{ID: 1, Name: "git"},
				{ID: 2, Name: "docker"},
				{ID: 3, Name: "kubernetes"},
			},
			expected: true,
		},
		{
			name: "same elements different order",
			skillsA: []models.Skill{
				{ID: 1, Name: "git"},
				{ID: 2, Name: "docker"},
			},
			skillsB: []models.Skill{
				{ID: 2, Name: "docker"},
				{ID: 1, Name: "git"},
			},
			expected: false,
		},
		{
			name:     "one empty one with elements",
			skillsA:  []models.Skill{},
			skillsB:  []models.Skill{{ID: 1, Name: "git"}},
			expected: false,
		},
		{
			name:     "zero values",
			skillsA:  []models.Skill{{}},
			skillsB:  []models.Skill{{ID: 0, Name: ""}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.SkillsEqual(tt.skillsA, tt.skillsB)
			if result != tt.expected {
				t.Errorf("SkillsEqual() = %v, expected %v", result, tt.expected)
				t.Errorf("skillsA: %+v", tt.skillsA)
				t.Errorf("skillsB: %+v", tt.skillsB)
			}
		})
	}
}