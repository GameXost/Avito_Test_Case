package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

// используется для создания пулл реквеста - получение тиммейтов по ID автора
func (t *TeamRepo) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	var team models.Team
	team.TeamName = name

	// запрос на получение данных для команды
	query := `SELECT u.user_id, u.username, u.is_active FROM teams AS t JOIN users AS u ON t.id = u.team_id WHERE t.name = $1`
	rows, err := t.pool.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("error in GetTeam %w", err)
	}
	defer rows.Close()

	// добавляем членов команды в слайс
	for rows.Next() {
		var member models.TeamMember
		err = rows.Scan(&member.UserId, &member.UserName, &member.IsActive)
		if err != nil {
			return nil, fmt.Errorf("error in GetTeam %w", err)
		}

		team.Members = append(team.Members, member)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in GetTeam %w", err)
	}

	// простая проверка на отсутствие команды в БД
	if len(team.Members) == 0 {
		return nil, models.ErrNotFound
	}

	return &team, nil

}
func (t *TeamRepo) AddTeam(ctx context.Context, team models.Team) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error in AddTeam %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	var teamID int
	query := `INSERT INTO teams (name) VALUES ($1) RETURNING id`
	err = tx.QueryRow(ctx, query, team.TeamName).Scan(&teamID)
	if err != nil {
		var pgErr *pgconn.PgError
		// находит первую ошибку в err, и сравнивает с кодом, что запись существует
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.ErrTeamExists
		}
		return fmt.Errorf("error in AddTeam %w", err)
	}
	// запихиваем пользователей в БД
	queryMember := `INSERT INTO users (user_id, username, is_active, team_id) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET username = excluded.username, is_active = excluded.is_active, team_id = excluded.team_id`
	for _, member := range team.Members {
		_, err = tx.Exec(ctx, queryMember, member.UserId, member.UserName, member.IsActive, teamID)
		if err != nil {
			return fmt.Errorf("error in AddTeam %w", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error in AddTeam %w", err)
	}
	return nil
}
