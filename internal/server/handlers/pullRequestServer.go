package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GameXost/Avito_Test_Case/internal/pkg/errHandle"
	"github.com/GameXost/Avito_Test_Case/internal/service"
	"github.com/GameXost/Avito_Test_Case/models"
	"net/http"
)

type PRHandler struct {
	prService *service.PullRequestService
}

type PRReq struct {
	PRID     string `json:"pull_request_id"`
	PRName   string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func NewPRHandler(prService *service.PullRequestService) *PRHandler {
	return &PRHandler{prService: prService}
}

func (pr *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req PRReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrInvalidJSON, "invalid JSON")
		return
	}
	if req.PRID == "" || req.PRName == "" || req.AuthorID == "" {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "required fields: id, name, autrho_id")
		return
	}

	pullrequest, err := pr.prService.CreatePR(r.Context(), req.PRID, req.PRName, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			errHandle.WriteError(w, http.StatusNotFound, models.ErrNotFound, "resource not found")
			return
		case errors.Is(err, models.ErrPRExists):
			errHandle.WriteError(w, http.StatusConflict, models.ErrPRExists, "PR id already exists")
			return
		default:
			errHandle.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(*pullrequest)
}

func (pr *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "invalid request body")
		return
	}

	if req.PullRequestID == "" {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "pull_request_id is required")
		return
	}

	pullRequest, err := pr.prService.Merge(r.Context(), req.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			errHandle.WriteError(w, http.StatusNotFound, models.ErrNotFound, "resource not found")
			return
		default:
			errHandle.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*pullRequest)
}

func (pr *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "invalid request body")
		return
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "prID and oldRev are required both")
		return
	}

	pullRequest, err := pr.prService.Reassign(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			errHandle.WriteError(w, http.StatusNotFound, models.ErrNotFound, "resource not found")
			return
		case errors.Is(err, models.ErrPRMerged):
			errHandle.WriteError(w, http.StatusConflict, models.ErrPRMerged, "cannot reassign on merged PR")
			return
		case errors.Is(err, models.ErrNotAssigned):
			errHandle.WriteError(w, http.StatusConflict, models.ErrNotAssigned, "reviewer is not assigned to this PR")
			return
		case errors.Is(err, models.ErrNoCandidate):
			errHandle.WriteError(w, http.StatusConflict, models.ErrNoCandidate, "no active replacement candidate in team")
			return
		default:
			errHandle.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*pullRequest)

}
