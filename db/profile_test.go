package db

import (
	"errors"
	"testing"
	"time"

	"github.com/flmailla/resume/models"
)

func TestGetProfiles(t *testing.T) {
	tests := []struct {
		name    string
		mockDB  *MockDB
		want    []models.Profile
		wantErr bool
	}{
		{
			name: "successful query with multiple profiles",
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
								*dest[1].(*string) = "Florent"
								*dest[2].(*string) = "Maillard"
								*dest[3].(*string) = "He"
								*dest[4].(*string) = "email@maillard.ch"
								*dest[5].(*string) = "Switzerland"
								*dest[6].(*int32) = int32(1000)
								*dest[7].(*string) = "Headline"
								*dest[8].(*string) = "About"
								*dest[9].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)

							} else if callCount == 2 {
								*dest[0].(*int64) = int64(2)
								*dest[1].(*string) = "Florence"
								*dest[2].(*string) = "Maillard"
								*dest[3].(*string) = "She"
								*dest[4].(*string) = "email2@maillard.ch"
								*dest[5].(*string) = "Switzerland"
								*dest[6].(*int32) = int32(1000)
								*dest[7].(*string) = "Headline"
								*dest[8].(*string) = "About"
								*dest[9].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)

							}
							return nil
						},
					}, nil
				},
			},
			want: []models.Profile{
				{
					ID:         1,
					FirstName:  "Florent",
					LastName:   "Maillard",
					Pronoun:    "He",
					Email:      "email@maillard.ch",
					Location:   "Switzerland",
					PostalCode: 1000,
					Headline:   "Headline",
					About:      "About",
					BirthDate:  time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
				},
				{
					ID:         2,
					FirstName:  "Florence",
					LastName:   "Maillard",
					Pronoun:    "She",
					Email:      "email2@maillard.ch",
					Location:   "Switzerland",
					PostalCode: 1000,
					Headline:   "Headline",
					About:      "About",
					BirthDate:  time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "database query error",
			mockDB: &MockDB{
				queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
					return nil, models.ErrDBRequestFailed
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
							return models.ErrScanFailed
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

			got, err := store.GetProfiles()

			if (err != nil) != tt.wantErr {
				t.Errorf("Error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Got %d profile, want %d profile", len(got), len(tt.want))
			}

			for i, profile := range got {
				if i < len(tt.want) {
					if profile != tt.want[i] {
						t.Errorf("Store.GetGetProfiles() got profile %+v, want %+v", profile, tt.want[i])
					}
				}
			}
		})
	}
}

func TestGetProfileById(t *testing.T) {
	tests := []struct {
		name    string
		mockDB  *MockDB
		want    models.Profile
		wantErr bool
	}{
		{
			name: "successful query profile",
			mockDB: &MockDB{
				queryRowFunc: func(query string, args ...interface{}) RowInterface {
					return &MockRow{
						scanFunc: func(dest ...interface{}) error {
							*dest[0].(*int64) = int64(1)
							*dest[1].(*string) = "Florent"
							*dest[2].(*string) = "Maillard"
							*dest[3].(*string) = "He"
							*dest[4].(*string) = "email3@maillard.ch"
							*dest[5].(*string) = "Switzerland"
							*dest[6].(*int32) = int32(1000)
							*dest[7].(*string) = "Headline"
							*dest[8].(*string) = "About"
							*dest[9].(*time.Time) = time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC)
							return nil
						},
					}
				},
			},
			want: models.Profile{
				ID:         1,
				FirstName:  "Florent",
				LastName:   "Maillard",
				Pronoun:    "He",
				Email:      "email3@maillard.ch",
				Location:   "Switzerland",
				PostalCode: 1000,
				Headline:   "Headline",
				About:      "About",
				BirthDate:  time.Date(2025, 1, 10, 23, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "no rows found error",
			mockDB: &MockDB{
				queryRowFunc: func(query string, args ...interface{}) RowInterface {
					return &MockRow{
						scanFunc: func(dest ...interface{}) error {
							return errors.New("No row")
						},
					}
				},
			},
			want:    models.Profile{},
			wantErr: true,
		},
		{
			name: "scan error",
			mockDB: &MockDB{
				queryRowFunc: func(query string, args ...interface{}) RowInterface {
					return &MockRow{
						scanFunc: func(dest ...interface{}) error {
							return errors.New("scan failed")
						},
					}
				},
			},
			want:    models.Profile{},
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
			want:    models.Profile{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(tt.mockDB)

			got, err := store.GetProfileById(1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("Profile.ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.FirstName != tt.want.FirstName {
				t.Errorf("Profile.FirstName = %v, want %v", got.FirstName, tt.want.FirstName)
			}
			if got.LastName != tt.want.LastName {
				t.Errorf("Profile.LastName = %v, want %v", got.LastName, tt.want.LastName)
			}
			if got.Pronoun != tt.want.Pronoun {
				t.Errorf("Profile.Pronoun = %v, want %v", got.Pronoun, tt.want.Pronoun)
			}
			if got.Email != tt.want.Email {
				t.Errorf("Profile.Email = %v, want %v", got.Email, tt.want.Email)
			}
			if got.Location != tt.want.Location {
				t.Errorf("Profile.Location = %v, want %v", got.Location, tt.want.Location)
			}
			if got.PostalCode != tt.want.PostalCode {
				t.Errorf("Profile.PostalCode = %v, want %v", got.PostalCode, tt.want.PostalCode)
			}
			if got.Headline != tt.want.Headline {
				t.Errorf("Profile.Headline = %v, want %v", got.Headline, tt.want.Headline)
			}
			if got.About != tt.want.About {
				t.Errorf("Profile.About = %v, want %v", got.About, tt.want.About)
			}
			if got.BirthDate != tt.want.BirthDate {
				t.Errorf("Profile.BirthDate = %v, want %v", got.BirthDate, tt.want.BirthDate)
			}

		})
	}
}
