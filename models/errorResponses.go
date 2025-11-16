package models

import "errors"

var (
	ErrTeamExists  = errors.New("TEAM_EXISTS")
	ErrPRExists    = errors.New("PR_EXISTS")
	ErrPRMerged    = errors.New("PR_VERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
	ErrNotFound    = errors.New("NOT_FOUND")
	ErrInvalidJSON = errors.New("INVALID_JSON")
	ErrValidation  = errors.New("VALIDATION_ERROR")
	ErrDefault     = errors.New("INTERNAL_ERROR")
)

type ErrorMessage struct {
	Code    error  `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}
