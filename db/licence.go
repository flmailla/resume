package db

import (
	"github.com/flmailla/resume/models"
)

func (s *Store) GetDistinctLicencesByProfile(profileId int) ([]models.Licence, error) {
	query := `SELECT DISTINCT l.id, l.title, l.issuer, l.issued_at, l.expires, l.licence_type
				FROM licence as l
				Where l.profile_id = ?`
	rows, err := s.db.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var licences []models.Licence

	for rows.Next() {
		var licence models.Licence
		if err := rows.Scan(&licence.ID,
			&licence.Title,
			&licence.Issuer,
			&licence.IssuedAt,
			&licence.Expires,
			&licence.LicenceType); err != nil {
			return licences, err
		}
		licences = append(licences, licence)
	}

	if err = rows.Err(); err != nil {
		return licences, err
	}
	return licences, nil
}
