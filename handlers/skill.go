package handlers

import (
	"net/http"
	"strconv"
	"github.com/flmailla/resume/models"
)

type SkillHandler struct {
	store storeHandler
}

func NewSkillHandler(store storeHandler) *SkillHandler {
	return &SkillHandler{store: store}
}

// @Summary Get all the skills 
// @Description Retrieve all the skills in the database
// @Tags Skills
// @Accept json
// @Produce json
// @Success 200 {object} models.Skill
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /skills [get]
func (h *SkillHandler) GetSkills(w http.ResponseWriter, r *http.Request) {
	profile, err := h.store.GetDistinctSkills()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrSkillsNotFetched.Error()})
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// @Summary Get a profile skills 
// @Description Retrieve all the skills for a given profile
// @Tags Skills
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Skill
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Param profile_id path int true "Profile ID"
// @Router /profiles/{profile_id}/skills [get]
func (h *SkillHandler) GetSkillsByProfile(w http.ResponseWriter, r *http.Request) {
	profileId, err := strconv.Atoi(r.PathValue("profile_id"))
	if err!= nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidId.Error()})
		return
	}
	skills, err := h.store.GetDistinctSkillsByProfile(profileId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrSkillsNotFetched.Error()})
		return
	}

	writeJSON(w, http.StatusOK, skills)
}

// @Summary Get the experience skills 
// @Description Retrieve all the skills for a given experience
// @Tags Skills
// @Tags Experience
// @Accept json
// @Produce json
// @Success 200 {object} models.Skill
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Param experience_id path int true "Experience ID"
// @Router /experience/{experience_id}/skills [get]
func (h *SkillHandler) GetSkillsByExperience(w http.ResponseWriter, r *http.Request) {
	experienceId, err := strconv.Atoi(r.PathValue("experience_id"))
	if err!= nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidId.Error()})
		return
	}
	experiences, err := h.store.GetDistinctSkillsByExperience(experienceId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrSkillsNotFetched.Error()})
		return
	}

	writeJSON(w, http.StatusOK, experiences)
}
