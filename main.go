package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/flmailla/resume/db"
	"github.com/flmailla/resume/handlers"
	"github.com/flmailla/resume/internal/auth"
	"github.com/flmailla/resume/logger"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		logger.Logger.Error("Received signal, shutting down...")
		os.Exit(0)
	}()

	if err := logger.InitLogger(); err != nil {
		panic("Failed to initialize logger")
	}

	if err := db.InitDB(); err != nil {
		logger.Logger.Error("Failed to initialize database", "error", err)
	}
	defer db.CloseDB()

	logger.Logger.Info("Application started")

	store := db.NewStoreFromSQLDB(db.DB)
	profileHandler := handlers.NewProfileHandler(store)
	experienceHandler := handlers.NewExperienceHandler(store)
	skillHandler := handlers.NewSkillHandler(store)
	educationHandler := handlers.NewEducationHandler(store)
	licenceHandler := handlers.NewLicenceHandler(store)
	healthHandler := handlers.NewHealthHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /profiles/{profile_id}", profileHandler.GetProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/experiences", experienceHandler.GetExperiencesByProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/skills", skillHandler.GetSkillsByProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/educations", educationHandler.GetEducationsByProfile)
	mux.HandleFunc("GET /profiles/{profile_id}/licences", licenceHandler.GetLicencesByProfile)
	mux.HandleFunc("GET /experiences/{experience_id}/skills", skillHandler.GetSkillsByExperience)
	mux.HandleFunc("GET /skills", skillHandler.GetSkills)
	mux.HandleFunc("GET /health", healthHandler.GetHealthStatus)

	validator := auth.NewJWTValidator("https://login.microsoftonline.com/df111d67-4cb1-4119-9f05-4c52e5e0e150/discovery/v2.0/keys")
	wrapped := validator.AuthMiddleware(mux)

	http.ListenAndServe("localhost:8090", wrapped)
}
