package db

import (
	"github.com/flmailla/resume/models"
)

func (s *Store) GetProfileById(ProfileId int) (*models.Profile, error) {
	var profile models.Profile
	query := "SELECT id, firstname, lastname, pronoun, email, location, postal_code, headline, about, birthdate FROM profile WHERE id = ?"
	err := s.db.QueryRow(query, ProfileId).Scan(
		&profile.ID, 
		&profile.FirstName, 
		&profile.LastName, 
		&profile.Pronoun,
		&profile.Email,
		&profile.Location,
		&profile.PostalCode,
		&profile.Headline,
		&profile.About,
		&profile.BirthDate)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *Store) GetProfiles() ([]models.Profile, error) {
	query := "SELECT id, firstname, lastname, pronoun, email, location, postal_code, headline, about, birthdate FROM profile"
	rows, err := s.db.Query(query)
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var profiles []models.Profile

	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(
			&profile.ID, 
			&profile.FirstName, 
			&profile.LastName, 
			&profile.Pronoun,
			&profile.Email,
			&profile.Location,
			&profile.PostalCode,
			&profile.Headline,
			&profile.About,
			&profile.BirthDate); err != nil {
				return profiles, err
			}
		profiles = append(profiles, profile)
	}
	 
	if err = rows.Err(); err != nil {
        return profiles, err
    }
    return profiles, nil
}