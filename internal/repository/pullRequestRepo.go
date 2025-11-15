package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/models"
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

func (pr *PullRequestRepo) CreatePR(ctx context.Context, request models.PullRequest) error {
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error in CreatePR %w", err)
	}
	defer tx.Rollback(ctx)
	queryPR := `INSERT INTO pull_request (pull_request_id, pull_request_name, author_id) VALUES ($1, $2, $3)`
	queryReviewer := `INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)`

	_, err = tx.Exec(ctx, queryPR, request.PullRequestID, request.PullRequestName, request.AuthorID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.ErrPRExists
		}
		return fmt.Errorf("error in CreatePR %w", err)
	}
	for _, reviewerID := range request.AssignedReviewers {
		_, err = tx.Exec(ctx, queryReviewer, request.PullRequestID, reviewerID)
		if err != nil {
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
	queryCheck := `SELECT EXISTS(SELECT 1 FROM pull_request WHERE pull_request_id = $1`
	err = tx.QueryRow(ctx, queryCheck, prID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
	}
	if !exists {
		return nil, models.ErrNotFound
	}

	querySet := `UPDATE pull_request SET status = 'MERGED', merged_at = now() WHERE pull_request_id = $1`
	_, err = tx.Exec(ctx, querySet, prID)
	if err != nil {
		return nil, fmt.Errorf("error in MergePR %w", err)
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

func (pr *PullRequestRepo) ReAssign(prID, reviewerNew string) (*models.PullRequest, error) {
	query := ``

}
