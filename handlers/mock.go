package handlers

import (
	"github.com/flmailla/resume/models"
)

type mockStore struct {
	GetDistinctEducationsByProfileFunc  func(profileId int) ([]models.Education, error)
	GetDistinctExperiencesByProfileFunc func(profileId int) ([]models.Experience, error)
	GetDistinctLicencesByProfileFunc    func(profileId int) ([]models.Licence, error)
	GetProfileFunc                      func(profileId int) (*models.Profile, error)
	GetDistinctSkillsFunc               func() ([]models.Skill, error)
	GetDistinctSkillsByProfileFunc      func(profileId int) ([]models.Skill, error)
	GetDistinctSkillsByExperienceFunc   func(experienceId int) ([]models.Skill, error)
}

func (m *mockStore) GetDistinctEducationsByProfile(profileId int) ([]models.Education, error) {
	if m.GetDistinctEducationsByProfileFunc != nil {
		return m.GetDistinctEducationsByProfileFunc(profileId)
	}
	return nil, models.ErrNotImplemented
}

func (m *mockStore) GetDistinctExperiencesByProfile(profileId int) ([]models.Experience, error) {
	if m.GetDistinctExperiencesByProfileFunc != nil {
		return m.GetDistinctExperiencesByProfileFunc(profileId)
	}
	return nil, models.ErrNotImplemented
}

func (m *mockStore) GetDistinctLicencesByProfile(profileId int) ([]models.Licence, error) {
	if m.GetDistinctLicencesByProfileFunc != nil {
		return m.GetDistinctLicencesByProfileFunc(profileId)
	}
	return nil, models.ErrNotImplemented
}

func (m *mockStore) GetProfileById(profileId int) (*models.Profile, error) {
	if m.GetProfileFunc != nil {
		return m.GetProfileFunc(profileId)
	}
	return &models.Profile{}, models.ErrNotImplemented
}

func (m *mockStore) GetDistinctSkills() ([]models.Skill, error) {
	if m.GetDistinctSkillsFunc != nil {
		return m.GetDistinctSkillsFunc()
	}
	return nil, models.ErrNotImplemented
}

func (m *mockStore) GetDistinctSkillsByProfile(profileId int) ([]models.Skill, error) {
	if m.GetDistinctSkillsByProfileFunc != nil {
		return m.GetDistinctSkillsByProfileFunc(profileId)
	}
	return nil, models.ErrNotImplemented
}

func (m *mockStore) GetDistinctSkillsByExperience(experienceId int) ([]models.Skill, error) {
	if m.GetDistinctSkillsByExperienceFunc != nil {
		return m.GetDistinctSkillsByExperienceFunc(experienceId)
	}
	return nil, models.ErrNotImplemented
}
