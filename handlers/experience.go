package handlers

import (
	"net/http"
	"strconv"
	"github.com/flmailla/resume/models"
	"github.com/flmailla/resume/logger"
)

type ExperienceHandler struct {
	store storeHandler
}

func NewExperienceHandler(store storeHandler) *ExperienceHandler {
	return &ExperienceHandler{store: store}
}

// @Summary Get a profile experiences 
// @Description Retrieve all the experiences for a given profile
// @Tags Experience
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Experience
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /profiles/{profile_id}/experiences [get]
// @Param profile_id path int true "Profile ID"
// @Security OAuth2Application
func (h *ExperienceHandler) GetExperiencesByProfile(w http.ResponseWriter, r *http.Request) {

	logger.Logger.Info("health endpoint requested")

	profileId, err := strconv.Atoi(r.PathValue("profile_id"))
	if err!= nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidId.Error()})
		logger.Logger.Warn("Experience endpoint", models.ErrInvalidId.Error(), profileId)
		return
	}
	profile, err := h.store.GetDistinctExperiencesByProfile(profileId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrExperiencesNotFetched.Error()})
		return
	}

	writeJSON(w, http.StatusOK, profile)
}