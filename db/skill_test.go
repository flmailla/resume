package db

import (
	"errors"
	"testing"
	"github.com/flmailla/resume/models"
)

func TestGetSkillsByXXX(t *testing.T) {
	tests := []struct {
		name    string
		mockDB  *MockDB
		want    []models.Skill
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
								*dest[1].(*string) = "Skill1"
							} else if callCount == 2 {
								*dest[0].(*int64) = int64(2)
								*dest[1].(*string) = "Skill2"
							}
							return nil
						},
					}, nil
				},
			},
			want: []models.Skill{
				{ID: 1, Name: "Skill1",},
				{ID: 2, Name: "Skill2",},
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
			
			got, err := store.GetDistinctSkillsByProfile(1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetSkillsByProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Store.GetSkillsByProfile() got %d skill, want %d skill", len(got), len(tt.want))
			}

			for i, skill := range got {
				if i < len(tt.want) {
					if skill != tt.want[i] {
						t.Errorf("Store.GetSkillsByProfile() got skill %+v, want %+v", skill, tt.want[i])
					}
				}
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(tt.mockDB)
			
			got, err := store.GetDistinctSkills()

			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetSkills() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Store.GetSkills() got %d skill, want %d skill", len(got), len(tt.want))
			}

			for i, skill := range got {
				if i < len(tt.want) {
					if skill != tt.want[i] {
						t.Errorf("Store.GetSkills() got skill %+v, want %+v", skill, tt.want[i])
					}
				}
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(tt.mockDB)
			
			got, err := store.GetDistinctSkillsByExperience(1)

			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetSkillsByExperience() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("Store.GetSkillsByExperience() got %d skill, want %d skill", len(got), len(tt.want))
			}

			for i, skill := range got {
				if i < len(tt.want) {
					if skill != tt.want[i] {
						t.Errorf("Store.GetSkillsByExperience() got skill %+v, want %+v", skill, tt.want[i])
					}
				}
			}
		})
	}
}

