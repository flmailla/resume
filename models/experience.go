package models

import (
	"time"
)

// Experience section in the resume
// Is linked to a profile
// Experiences lead to skills
type Experience struct {
	ID          int64
	Title       string
	Company     string
	StartDate   time.Time
	EndDate     time.Time
	Location    string
	Description string
	Skills      []Skill
	Profile     Profile
}
