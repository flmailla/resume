package handlers

import (
	"net/http"
	"strconv"

	"github.com/flmailla/resume/logger"
	"github.com/flmailla/resume/models"
)

type LicenceHandler struct {
	store storeHandler
}

func NewLicenceHandler(store storeHandler) *LicenceHandler {
	return &LicenceHandler{store: store}
}

// @Summary Get a profile Licences
// @Description Retrieve all the Licences for a given profile
// @Tags Licence
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Licence
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Param profile_id path int true "Profile ID"
// @Router /profiles/{profile_id}/licences [get]
// @Security OAuth2Application
func (h *LicenceHandler) GetLicencesByProfile(w http.ResponseWriter, r *http.Request) {
	profileId, err := strconv.Atoi(r.PathValue("profile_id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidId.Error()})
		logger.Logger.Warn("Licence endpoint", models.ErrInvalidId.Error(), profileId)
		return
	}
	licences, err := h.store.GetDistinctLicencesByProfile(profileId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": models.ErrLicencesNotFetched.Error()})
		return
	}

	writeJSON(w, http.StatusOK, licences)
}
