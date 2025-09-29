package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/flmailla/resume/logger"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Logger.Error("Failed to marshal JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(status)
	w.Write(data)
}