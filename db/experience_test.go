package db

import (
	"errors"
	"time"
	"testing"
	"github.com/flmailla/resume/models"
)

func TestGetExperiencesByProfile(t *testing.T) {
	tests := []struct {
		name    string
		mockDB  *MockDB
		want    []models.Experience
		wantErr bool
	}{
		{
			name: "successful query with multiple experiences",
			mockDB: &MockDB{
				queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
					callCount := 0
					return &MockRows{
						nextFunc: func() bool {
							callCount++
							return callCount <= 2 
						},
						scanFunc: func(dest ...interface{}) error {
							if callCount == 1 {
								*dest[0].(*int64) = int64(1)
								*dest[1].(*string) = "Job1"
								*dest[2].(*string) = "Company1"
								*dest[3].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[4].(*time.Time) = time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[5].(*string) = "Lausanne"
								*dest[6].(*string) = "Just a job"
							} else if callCount == 2 {
								*dest[0].(*int64) = int64(2)
								*dest[1].(*string) = "Job2"
								*dest[2].(*string) = "Company2"
								*dest[3].(*time.Time) = time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[4].(*time.Time) = time.Time{}
								*dest[5].(*string) = "Sion"
								*dest[6].(*string) = "Just another job"
							}
							return nil
						},
					}, nil
				},
			},
			want: []models.Experience{
				{
					ID: 1, 
					Title: "Job1", 
					Company: "Company1", 
					StartDate: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					EndDate: time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC), 
					Location: "Lausanne",
					Description: "Just a job",
					Skills: []models.Skill{}},
				{
					ID: 2, 
					Title: "Job2", 
					Company: "Company2", 
					StartDate: time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC), 
					EndDate: time.Time{},
					Location: "Sion",
					Description: "Just another job",
					Skills: []models.Skill{}},
			},
			wantErr: false,
		},
		{
			name: "database query error",
			mockDB: &MockDB{
				queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
					return nil, errors.New("database connection failed")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "scan error",
			mockDB: &MockDB{
				queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
					return &MockRows{
						nextFunc: func() bool {
							return true
						},
						scanFunc: func(dest ...interface{}) error {
							return errors.New("scan failed")
						},
					}, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DB error",
			mockDB: &MockDB{
				queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
					return &MockRows{
						nextFunc: func() bool {
							return false
						},
						scanFunc: func(dest ...interface{}) error {
							return nil
						},
						errFunc: func() error {
							return models.ErrDBRequestFailed
						},
					}, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(tt.mockDB)
			
			got, err := store.GetDistinctExperiencesByProfile(1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Got %d experience, want %d experience", len(got), len(tt.want))
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
		})
	}
}