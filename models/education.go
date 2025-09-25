package models

import (
	"time"
)


// Education section in the resume
// Is linked to a profile
type Education struct {
	ID              int64
	Title           string
	Issued          time.Time
	Description     string
	Profile         Profile
}