package models

import (
	"time"
)

type LicenceType string

const (
	LICENCE            LicenceType = "Licence"
	CERTIFICATION      LicenceType = "Certification" 
)

// Licence section in the resume
// Is linked to a profile
type Licence struct {
	ID              int64
	Title           string
	Issuer          string
	IssuedAt        time.Time
	Expires         time.Time
    LicenceType     LicenceType
	Profile         Profile
}