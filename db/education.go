package db

import (
	"github.com/flmailla/resume/models"
)

func (s *Store) GetDistinctEducationsByProfile(profileId int) ([]models.Education, error) {
	query := `SELECT DISTINCT e.id, e.title, e.issued_at, e.description
				FROM education as e
				Where e.profile_id = ?`
	rows, err := s.db.Query(query, profileId)
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var educations []models.Education

	for rows.Next() {
		var education models.Education
		if err := rows.Scan(&education.ID,
							&education.Title,
							&education.Issued,
							&education.Description); err != nil {
				return educations, err
			}
		educations = append(educations, education)
	}
	 
	if err = rows.Err(); err != nil {
        return educations, err
    }
    return educations, nil
}