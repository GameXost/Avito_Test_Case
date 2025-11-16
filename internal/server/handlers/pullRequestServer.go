package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GameXost/Avito_Test_Case/internal/server"
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
		server.WriteError(w, http.StatusBadRequest, models.ErrInvalidJSON, "invalid JSON")
		return
	}
	if req.PRID == "" || req.PRName == "" || req.AuthorID == "" {
		server.WriteError(w, http.StatusBadRequest, models.ErrValidation, "required fields: id, name, autrho_id")
		return
	}

	pullrequest, err := pr.prService.CreatePR(r.Context(), req.PRID, req.PRName, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			server.WriteError(w, http.StatusNotFound, models.ErrNotFound, "resource not found")
			return
		case errors.Is(err, models.ErrPRExists):
			server.WriteError(w, http.StatusConflict, models.ErrPRExists, "PR id already exists")
			return
		default:
			server.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(*pullrequest)
}

func (pr *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	prID := r.URL.Query().Get("pull_request_id")
	if prID == "" {
		server.WriteError(w, http.StatusBadRequest, models.ErrValidation, "pull_request_id is required")
		return
	}
	pullRequest, err := pr.prService.Merge(r.Context(), prID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			server.WriteError(w, http.StatusNotFound, models.ErrNotFound, "resource not found")
			return
		default:
			server.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*pullRequest)
}

func (pr *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	prID := r.URL.Query().Get("pull_request_id")
	oldRev := r.URL.Query().Get("old_reviewer_id")
	if prID == "" || oldRev == "" {
		server.WriteError(w, http.StatusBadRequest, models.ErrValidation, "prID and oldRev are required both")
		return
	}
	pullRequest, err := pr.prService.Reassign(r.Context(), prID, oldRev)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			server.WriteError(w, http.StatusNotFound, models.ErrNotFound, "resource not found")
			return
		case errors.Is(err, models.ErrPRMerged):
			server.WriteError(w, http.StatusConflict, models.ErrPRMerged, "cannot reassign on merged PR")
			return
		case errors.Is(err, models.ErrNotAssigned):
			server.WriteError(w, http.StatusConflict, models.ErrNotAssigned, "reviewer is not assigned to this PR")
			return
		case errors.Is(err, models.ErrNoCandidate):
			server.WriteError(w, http.StatusConflict, models.ErrNoCandidate, "no active replacement candidate in team")
			return
		default:
			server.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*pullRequest)

}
