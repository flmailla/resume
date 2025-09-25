package handlers

import (
	"net/http"
	"strconv"
	"github.com/flmailla/resume/models"
)

type ProfileHandler struct {
	store storeHandler
}

func NewProfileHandler(store storeHandler) *ProfileHandler {
	return &ProfileHandler{store: store}
}

// @Summary Get a profile 
// @Description Retrieve the profil information
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Profile
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Param profile_id path int true "Profile ID"
// @Router /profiles/{profile_id} [get]
// @Security OAuth2Application
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    profileId, err := strconv.Atoi(r.PathValue("profile_id"))
	if err!= nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidId.Error()})
		return
	}
	profile, err := h.store.GetProfileById(profileId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrProfileNotFetched.Error()})
		return
	}

	writeJSON(w, http.StatusOK, profile)
}