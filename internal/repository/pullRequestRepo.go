package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepo struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepo(pool *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{
		pool: pool,
	}
}

// убрать часть с добавлением в БД ревьюеров, использовать AssignReviewer в сервисе
func (pr *PullRequestRepo) CreatePR(ctx context.Context, request models.PullRequest) error {
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error in CreatePR %w", err)
	}
	defer tx.Rollback(ctx)
	queryPR := `INSERT INTO pull_request (pull_request_id, pull_request_name, author_id) VALUES ($1, $2, $3)`

	_, err = tx.Exec(ctx, queryPR, request.PullRequestID, request.PullRequestName, request.AuthorID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.ErrPRExists
		}
		return fmt.Errorf("error in CreatePR %w", err)
	}

	for _, reviewer := range request.AssignedReviewers {
		err = pr.AssignReviewer(ctx, tx, request.PullRequestID, reviewer)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return models.ErrNotFound
			}
			return fmt.Errorf("error in CreatePR %w", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error in CreatePR %w", err)
	}
	return nil
}

func (pr *PullRequestRepo) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pullRequest models.PullRequest
	pullRequest.PullRequestID = prID
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	queryCheck := `SELECT EXISTS(SELECT 1 FROM pull_request WHERE pull_request_id = $1)`
	err = tx.QueryRow(ctx, queryCheck, prID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}
	if !exists {
		return nil, models.ErrNotFound
	}

	var status string
	queryIsOpen := `SELECT status FROM pull_request WHERE pull_request_id = $1`
	err = tx.QueryRow(ctx, queryIsOpen, prID).Scan(&status)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}

	if status == "OPEN" {
		querySet := `UPDATE pull_request SET status = 'MERGED', merged_at = now() WHERE pull_request_id = $1`
		_, err = tx.Exec(ctx, querySet, prID)
		if err != nil {
			return nil, fmt.Errorf("error in MergePR %w", err)
		}
	}
	queryGetPR := `SELECT pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at FROM pull_request AS pr WHERE pr.pull_request_id = $1`
	err = tx.QueryRow(ctx, queryGetPR, prID).Scan(&pullRequest.PullRequestName, &pullRequest.AuthorID, &pullRequest.Status, &pullRequest.CreatedAt, &pullRequest.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}

	queryGetRev := `SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1`
	rows, err := tx.Query(ctx, queryGetRev, prID)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var reviewerID string
		err = rows.Scan(&reviewerID)
		if err != nil {
			return nil, fmt.Errorf("error in MergePR %w", err)
		}
		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewerID)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}
	return &pullRequest, nil
}

// проверку на отсутсвие
func (pr *PullRequestRepo) GetTeamNameByUserID(ctx context.Context, userID string) (string, error) {
	query := `SELECT t.name FROM users AS u JOIN teams AS t ON u.team_id = t.id WHERE u.user_id = $1`
	var teamName string
	err := pr.pool.QueryRow(ctx, query, userID).Scan(&teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrNotFound
		}
		return "", fmt.Errorf("error in GetTeamIDByUserID %w", err)
	}
	return teamName, nil
}

func (pr *PullRequestRepo) GetPRInfo(ctx context.Context, prID string) (*models.PullRequest, error) {
	query := `SELECT pr.pull_request_name, pr.author_id, pr.status FROM pull_request AS pr WHERE pr.pull_request_id = $1`
	var pullReq models.PullRequest
	err := pr.pool.QueryRow(ctx, query, prID).Scan(&pullReq.PullRequestName, &pullReq.AuthorID, &pullReq.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("error in GetPRInfo %w", err)
	}
	return &pullReq, nil
}

func (pr *PullRequestRepo) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	query := `SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1`
	rows, err := pr.pool.Query(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("error in GetReviewers %w", err)
	}
	defer rows.Close()
	var reviewers []string
	for rows.Next() {
		var reviewer string
		err = rows.Scan(&reviewer)
		if err != nil {
			return nil, fmt.Errorf("error in GetReviewers %w", err)
		}
		reviewers = append(reviewers, reviewer)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error in GetReviewers %w", err)
	}
	return reviewers, nil
}

func (pr *PullRequestRepo) DeleteReviewer(ctx context.Context, tx pgx.Tx, revID, prID string) error {
	queryDel := `DELETE FROM pr_reviewers WHERE reviewer_id = $1 AND pr_id = $2`
	res, err := tx.Exec(ctx, queryDel, revID, prID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return models.ErrNotAssigned
		}
		return fmt.Errorf("error in DELETE reviewer %w", err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotAssigned
	}

	return nil
}

func (pr *PullRequestRepo) AssignReviewer(ctx context.Context, tx pgx.Tx, prID, revNew string) error {
	queryIns := `INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES($1, $2)`
	res, err := tx.Exec(ctx, queryIns, prID, revNew)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return models.ErrNotFound
			}
		}
		return fmt.Errorf("error in AssignReviewer %w", err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotFound // PR не существует
	}
	return nil
}

func (pr *PullRequestRepo) IsMerged(ctx context.Context, prID string) (string, error) {
	var status string
	query := `SELECT status FROM pull_request WHERE pull_request_id = $1`
	err := pr.pool.QueryRow(ctx, query, prID).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("error in IsMerged %w", err)
	}
	return status, nil
}

func (pr *PullRequestRepo) DelAndAssign(ctx context.Context, oldRevID, prID, newRevID string) error {
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error in DelAndAssign %w", err)
	}
	defer tx.Rollback(ctx)
	err = pr.DeleteReviewer(ctx, tx, oldRevID, prID)
	if err != nil {
		if errors.Is(err, models.ErrNotAssigned) {
			return models.ErrNotAssigned
		}
		return fmt.Errorf("error in DelAndAssign %w", err)
	}
	err = pr.AssignReviewer(ctx, tx, prID, newRevID)
	if err != nil {
		return fmt.Errorf("error in DelAndAssign %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error in DelAndAssign %w", err)
	}

	return nil
}
