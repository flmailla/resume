package db

import (
	"errors"
	"time"
	"testing"
	"github.com/flmailla/resume/models"
)

func TestGetEducationsByProfile(t *testing.T) {
	tests := []struct {
		name    string
		mockDB  *MockDB
		want    []models.Education
		wantErr bool
	}{
		{
			name: "successful query with multiple educations",
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
								*dest[1].(*string) = "University"
								*dest[2].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[3].(*string) = "A university journey"
							} else if callCount == 2 {
								*dest[0].(*int64) = int64(2)
								*dest[1].(*string) = "University again"
								*dest[2].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[3].(*string) = "Another university journey"
							}
							return nil
						},
					}, nil
				},
			},
			want: []models.Education{
				{ID: 1, Title: "University", Issued: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), Description: "A university journey"},
				{ID: 2, Title: "University again", Issued: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), Description: "Another university journey"},
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
		},{
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
			
			got, err := store.GetDistinctEducationsByProfile(1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Store.GetUsers() got %d users, want %d education", len(got), len(tt.want))
			}

			for i, education := range got {
				if i < len(tt.want) {
					if education != tt.want[i] {
						t.Errorf("Store.GetDistinctEducationsByProfile() got education %+v, want %+v", education, tt.want[i])
					}
				}
			}
		})
	}
}