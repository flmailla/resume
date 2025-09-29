package db

import (
	"github.com/flmailla/resume/models"
)

func (s *Store) GetDistinctSkills() ([]models.Skill, error) {
	query := "SELECT DISTINCT id, name FROM skill"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill

	for rows.Next() {
		var skill models.Skill
		if err := rows.Scan(
			&skill.ID,
			&skill.Name); err != nil {
			return skills, err
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return skills, err
	}
	return skills, nil
}

func (s *Store) GetDistinctSkillsByProfile(profileId int) ([]models.Skill, error) {
	query := `SELECT DISTINCT s.id, s.name FROM skill as s
				JOIN skill_experience as se ON se.skill_id = s.id
				JOIN experience as e ON e.id = se.experience_id
				Where e.profile_id = ?`
	rows, err := s.db.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill

	for rows.Next() {
		var skill models.Skill
		if err := rows.Scan(&skill.ID, &skill.Name); err != nil {
			return skills, err
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return skills, err
	}
	return skills, nil
}

func (s *Store) GetDistinctSkillsByExperience(experienceId int) ([]models.Skill, error) {
	query := `SELECT DISTINCT s.id, s.name FROM skill as s
				JOIN skill_experience as se ON se.skill_id = s.id
				JOIN experience as e ON e.id = se.experience_id
				Where e.id = ?`
	rows, err := s.db.Query(query, experienceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill

	for rows.Next() {
		var skill models.Skill
		if err := rows.Scan(&skill.ID, &skill.Name); err != nil {
			return skills, err
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return skills, err
	}
	return skills, nil
}
