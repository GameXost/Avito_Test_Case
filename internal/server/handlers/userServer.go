package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GameXost/Avito_Test_Case/internal/pkg/errHandle"
	"github.com/GameXost/Avito_Test_Case/internal/service"
	"github.com/GameXost/Avito_Test_Case/models"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}
type UserReq struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
type UserResp struct {
	UserID       string                    `json:"user_id"`
	PullRequests []models.PullRequestShort `json:"pull_requests"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req UserReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "invalid JSON")
		return
	}

	if req.UserID == "" {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "user_id is required")
		return
	}

	user, err := u.userService.SetActive(r.Context(), req.UserID, req.IsActive)
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
	_ = json.NewEncoder(w).Encode(*user)
}

func (u *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		errHandle.WriteError(w, http.StatusBadRequest, models.ErrValidation, "user_id is required")
		return
	}

	pullRequests, err := u.userService.GetUserReviews(r.Context(), userID)
	if err != nil {
		errHandle.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
		return
	}

	resp := UserResp{UserID: userID, PullRequests: pullRequests}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)

}
