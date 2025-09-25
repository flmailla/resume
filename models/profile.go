package models

import (
	"time"
)

// Basically, me
type Profile struct {
	ID             int64
	FirstName      string
	LastName       string
	BirthDate      time.Time
	Pronoun        string
	Email          string
	Location       string
	PostalCode     int32
	Headline       string
	About          string
}

