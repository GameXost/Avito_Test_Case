package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/internal/repository"
	"github.com/GameXost/Avito_Test_Case/models"
)

type TeamService struct {
	teamRepo repository.TeamRepo
}

func NewTeamService(teamRepo repository.TeamRepo) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (t *TeamService) CreateTeam(ctx context.Context, teamName string, teamMembers []models.TeamMember) (*models.Team, error) {
	team := models.Team{TeamName: teamName, Members: teamMembers}
	err := t.teamRepo.AddTeam(ctx, team)
	// тут наверное проверка на то, что именно вернулось нужна, если наша ошибка, то определенный статус код отправляем или нет, т.к на уровне репозитория тут уже всё учли
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTeamExists):
			return nil, models.ErrTeamExists
		default:
			return nil, fmt.Errorf("Create team: %w", err)
		}
	}
	return &team, nil
}

func (t *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := t.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrNotFound
		default:
			return nil, fmt.Errorf("GetTeam: %w", err)
		}
	}
	return team, nil
}
