package db

import (
	"github.com/flmailla/resume/models"
)

func (s *Store) GetDistinctExperiencesByProfile(profileId int) ([]models.Experience, error) {
	query := `SELECT DISTINCT e.id, e.title, e.company, e.start_date, e.end_date, e.location, e.description
				FROM experience as e
				Where e.profile_id = ?`
	rows, err := s.db.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var experiences []models.Experience

	for rows.Next() {
		var experience models.Experience
		if err := rows.Scan(&experience.ID,
			&experience.Title,
			&experience.Company,
			&experience.StartDate,
			&experience.EndDate,
			&experience.Location,
			&experience.Description); err != nil {
			return experiences, err
		}
		experiences = append(experiences, experience)
	}

	if err = rows.Err(); err != nil {
		return experiences, err
	}
	return experiences, nil
}
