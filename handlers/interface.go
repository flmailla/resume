package handlers

import (
	"github.com/flmailla/resume/models"
)

type storeHandler interface {
	GetDistinctEducationsByProfile(profileId int) ([]models.Education, error)
	GetDistinctExperiencesByProfile(profileId int) ([]models.Experience, error)
	GetDistinctLicencesByProfile(profileId int) ([]models.Licence, error)
	GetProfileById(profileId int) (*models.Profile, error)
	GetDistinctSkills() ([]models.Skill, error)
	GetDistinctSkillsByProfile(profileId int) ([]models.Skill, error)
	GetDistinctSkillsByExperience(experienceId int) ([]models.Skill, error)
}
