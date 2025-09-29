package models

import "reflect"

// Skill represents a skill in the system
// @Description Skill information
// @name Skill
type Skill struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"git"`
}

func SkillsEqual(a, b []Skill) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}
