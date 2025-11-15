package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (u *UserRepo) GetPR(ctx context.Context, userID string) (string, []models.PullRequestShort, error) {
	var PRequests []models.PullRequestShort
	// тело запроса: данные для ПуллРеквеста
	query := `
	SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
	FROM pr_reviewers AS rev
	JOIN pull_request AS pr ON rev.pr_id = pr.pull_request_id
	WHERE rev.reviewer_id = $1`

	rows, err := u.pool.Query(ctx, query, userID)
	if err != nil {
		return "", nil, fmt.Errorf("error in GetPR %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var PRShort models.PullRequestShort
		err = rows.Scan(&PRShort.PullRequestID, &PRShort.PullRequestName, &PRShort.AuthorID, &PRShort.Status)
		if err != nil {
			return "", nil, fmt.Errorf("error in GetPR %w", err)
		}
		PRequests = append(PRequests, PRShort)
	}
	err = rows.Err()
	if err != nil {
		return "", nil, fmt.Errorf("error in GetPR %w", err)
	}
	return userID, PRequests, nil
}

func (u *UserRepo) SetActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in SetActive %w", err)
	}
	defer tx.Rollback(ctx)

	var user models.User
	query := `UPDATE users set is_active = $1 WHERE user_id = $2`
	_, err = tx.Exec(ctx, query, isActive, userID)
	if err != nil {
		return nil, fmt.Errorf("error in SetActive %w", err)
	}

	query = `SELECT u.user_id, u.username, t.name, u.is_active FROM users AS u JOIN teams AS t ON u.team_id = t.id WHERE u.user_id = $1`
	err = tx.QueryRow(ctx, query, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		// проверка на случай, если пользователь не найден
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("error in SetActive %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in SetActive %w", err)
	}
	return &user, nil
}
