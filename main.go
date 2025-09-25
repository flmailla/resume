package main

import (
	"log"
	"net/http"
	"github.com/flmailla/resume/db"
	"github.com/flmailla/resume/handlers"
	"github.com/flmailla/resume/internal/auth"
)

// @title Resume API
// @version 1.0.0
// @description Displays my resume main sections as APIs
// @host localhost:8080
// @BasePath /resume/v1

// @contact.name flmailla
// @contact.email florent@maillard.icu

// @securityDefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	store := db.NewStoreFromSQLDB(db.DB)
	profileHandler := handlers.NewProfileHandler(store)
	experienceHandler := handlers.NewExperienceHandler(store)
	skillHandler := handlers.NewSkillHandler(store)
	educationHandler := handlers.NewEducationHandler(store)
	licenceHandler := handlers.NewLicenceHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /profiles/{profile_id}", profileHandler.GetProfile)
    mux.HandleFunc("GET /profiles/{profile_id}/experiences", experienceHandler.GetExperiencesByProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/skills", skillHandler.GetSkillsByProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/educations", educationHandler.GetEducationsByProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/licences", licenceHandler.GetLicencesByProfile)
	mux.HandleFunc("GET /experiences/{experience_id}/skills", skillHandler.GetSkillsByExperience)
	mux.HandleFunc("GET /skills", skillHandler.GetSkills)

	validator := auth.NewJWTValidator("https://login.microsoftonline.com/df111d67-4cb1-4119-9f05-4c52e5e0e150/discovery/v2.0/keys")
	wrapped := validator.AuthMiddleware(mux)

	http.ListenAndServe("localhost:8090", wrapped)
}