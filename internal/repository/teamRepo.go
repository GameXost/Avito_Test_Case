package repository

import (
	"context"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

func (t *TeamRepo) GetTeam(ctx context.Context, name string) (models.Team, error) {
	var res models.Team
	res.TeamName = name
	query := `SELECT users.user_id, users.username, users.is_active FROM teams JOIN users ON teams.id = users.team_id WHERE teams.name = $1`
	rows, err := t.pool.Query(ctx, query, name)
	if err != nil {
		return models.Team{}, fmt.Errorf("error in GetTeam %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var member models.TeamMember
		err = rows.Scan(&member.UserId, &member.UserName, &member.IsActive)
		if err != nil {
			return models.Team{}, fmt.Errorf("error in GetTeam %w", err)
		}

		res.Members = append(res.Members, member)
	}
	if err = rows.Err(); err != nil {
		return models.Team{}, fmt.Errorf("error in GetTeam %w", err)
	}

	return res, nil

}
func (t *TeamRepo) AddTeam(ctx context.Context, team models.Team) (models.Team, error) {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return models.Team{}, fmt.Errorf("error in AddTeam %w", err)
	}
	defer tx.Rollback(ctx)
	var id int
	query := `INSERT INTO teams (name) VALUES ($1) RETURNING id`
	err = tx.QueryRow(ctx, query, team.TeamName).Scan(&id)
	if err != nil {
		return models.Team{}, fmt.Errorf("error in AddTeam %w", err)
	}
	queryMember := `INSERT INTO users (user_id, username, is_active, team_id) VALUES ($1, $2, $3, $4)`
	for _, member := range team.Members {
		_, err = tx.Exec(ctx, queryMember, member.UserId, member.UserName, member.IsActive, id)
		if err != nil {
			return models.Team{}, fmt.Errorf("error in AddTeam %w", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return models.Team{}, fmt.Errorf("error in AddTeam %w", err)
	}
	return team, nil
}
