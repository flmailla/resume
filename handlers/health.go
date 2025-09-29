package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/flmailla/resume/logger"
)

type HealthHandler struct {
	store storeHandler
}

type responseStatus struct {
	Status	    string
}

func NewHealthHandler(store storeHandler) *HealthHandler {
	return &HealthHandler{store: store}
}

// @Summary Get a status about the service
// @Description Get the health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.Licence
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {

	logger.Logger.Info("health endpoint requested")

	response := map[string]string{
        "Status": "healthy",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
	