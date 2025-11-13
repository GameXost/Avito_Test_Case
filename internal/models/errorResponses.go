package models

import "errors"

var (
	ErrTeamExists  = errors.New("TEAM_EXISTS")
	ErrPRExists    = errors.New("PR_EXISTS")
	ErrPRMerged    = errors.New("PR_VERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
	ErrNotFound    = errors.New("NOT_FOUND")
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
