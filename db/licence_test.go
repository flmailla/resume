package db

import (
	"errors"
	"time"
	"testing"
	"github.com/flmailla/resume/models"
)

func TestGetLicencesByProfile(t *testing.T) {
	tests := []struct {
		name    string
		mockDB  *MockDB
		want    []models.Licence
		wantErr bool
	}{
		{
			name: "successful query with multiple licences",
			mockDB: &MockDB{
				queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
					callCount := 0
					return &MockRows{
						nextFunc: func() bool {
							callCount++
							return callCount <= 3 
						},
						scanFunc: func(dest ...interface{}) error {
							if callCount == 1 {
								*dest[0].(*int64) = int64(1)
								*dest[1].(*string) = "Job1"
								*dest[2].(*string) = "Company1"
								*dest[3].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[4].(*time.Time) = time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[5].(*models.LicenceType) = models.LICENCE
							} else if callCount == 2 {
								*dest[0].(*int64) = int64(2)
								*dest[1].(*string) = "Job2"
								*dest[2].(*string) = "Company2"
								*dest[3].(*time.Time) = time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[4].(*time.Time) = time.Time{}
								*dest[5].(*models.LicenceType) = models.CERTIFICATION
							} else if callCount == 3 {
								*dest[0].(*int64) = int64(3)
								*dest[1].(*string) = "Job3"
								*dest[2].(*string) = "Company3"
								*dest[3].(*time.Time) = time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC)
								*dest[4].(*time.Time) = time.Time{}
								*dest[5].(*models.LicenceType) = ""
							}
							return nil
						},
					}, nil
				},
			},
			want: []models.Licence{
				{
					ID: 1, 
					Title: "Job1", 
					Issuer: "Company1", 
					IssuedAt: time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC), 
					Expires: time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC), 
					LicenceType: models.LICENCE,
				},
				{
					ID: 2, 
					Title: "Job2", 
					Issuer: "Company2", 
					IssuedAt: time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC), 
					Expires: time.Time{},
					LicenceType: models.CERTIFICATION,
				},
				{
					ID: 3, 
					Title: "Job3", 
					Issuer: "Company3", 
					IssuedAt: time.Date(2026, 1, 10, 23, 0, 0, 0, time.UTC), 
					Expires: time.Time{},
					LicenceType: "",
				},
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
			
			got, err := store.GetDistinctLicencesByProfile(1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Got %d licence, want %d licence", len(got), len(tt.want))
			}

			for i, licence := range got {
				if i < len(tt.want) {
					if licence != tt.want[i] {
						t.Errorf("Store.GetDistinctLicencesByProfile() got licence %+v, want %+v", licence, tt.want[i])
					}
				}
			}
		})
	}
}