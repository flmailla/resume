package models

import "errors"

var (
	ErrNotImplemented        = errors.New("not implemented")
	ErrProfileNotFound       = errors.New("profile not found")
	ErrSkillsNotFound        = errors.New("no skill found")
	ErrSkillsNotFetched      = errors.New("fail to fetch skills")
	ErrEducationsNotFetched  = errors.New("failed to fetch educations")
	ErrProfileNotFetched     = errors.New("failed to fetch profile")
	ErrExperiencesNotFetched = errors.New("failed to fetch experiences")
	ErrLicencesNotFetched    = errors.New("failed to fetch licences")
	ErrUnknown               = errors.New("unknown error")
	ErrUnmarshal             = errors.New("failed to unmarshal response body: %v")
	ErrInvalidId             = errors.New("invalid id")
	ErrDBRequestFailed       = errors.New("db request failed")
	ErrScanFailed            = errors.New("db scan failed")
	ErrNoTokenSent           = errors.New("no token sent")
	ErrNotBearer             = errors.New("not a bearer token")
	ErrUnauthorized          = errors.New("unauthorized")
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"User not found"`
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"The requested profile could not be found"`
} // @name ErrorResponsNotFound

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
} // @name SuccessResponse
