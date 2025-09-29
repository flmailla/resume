package handlers

import (
	"net/http"
	"strconv"
	"github.com/flmailla/resume/models"
	"github.com/flmailla/resume/logger"
)

type EducationHandler struct {
	store storeHandler
}

func NewEducationHandler(store storeHandler) *EducationHandler {
	return &EducationHandler{store: store}
}

// @Summary Get a profile educations 
// @Description Retrieve all the education lines for a given profile
// @Tags Education
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Education
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /profiles/{profile_id}/educations [get]
// @Param profile_id path int true "Profile ID"
// @Security OAuth2Application
func (h *EducationHandler) GetEducationsByProfile(w http.ResponseWriter, r *http.Request) {
	profileId, err := strconv.Atoi(r.PathValue("profile_id"))
	if err != nil {
		logger.Logger.Error(err.Error())
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidId.Error()})
		return
	}
	educations, err := h.store.GetDistinctEducationsByProfile(profileId)
	if err != nil {
		logger.Logger.Error(err.Error())
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrEducationsNotFetched.Error(), "detail": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, educations)
}