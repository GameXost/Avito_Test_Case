package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GameXost/Avito_Test_Case/internal/server"
	"github.com/GameXost/Avito_Test_Case/internal/service"
	"github.com/GameXost/Avito_Test_Case/models"
	"net/http"
)

type TeamHandler struct {
	teamService *service.TeamService
}

type TeamReq struct {
	TeamName string              `json:"team_name"`
	Members  []models.TeamMember `json:"members"`
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (t *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req TeamReq

	// Ошибка JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteError(w, http.StatusBadRequest, models.ErrInvalidJSON, "invalid JSON")
		return
	}

	// Валидация также шлём 400 ошибку
	if req.TeamName == "" || req.Members == nil {
		server.WriteError(w, http.StatusBadRequest, models.ErrValidation, "team_name is required")
		return
	}

	team, err := t.teamService.CreateTeam(r.Context(), req.TeamName, req.Members)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTeamExists):
			server.WriteError(w, http.StatusBadRequest, models.ErrTeamExists, "team already exists")
			return
		// на любую другую ошибку шлем 500 респонс
		default:
			server.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}

	// OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(team)
}

func (t *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		server.WriteError(w, http.StatusBadRequest, models.ErrValidation, "team_name is required")
		return
	}

	team, err := t.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			server.WriteError(w, http.StatusNotFound, models.ErrNotFound, "team not found")
			return
		default:
			server.WriteError(w, http.StatusInternalServerError, models.ErrDefault, "unexpected server error")
			return
		}
	}
	res := TeamReq{TeamName: teamName, Members: team.Members}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
