package repository

import (
	"context"
	"github.com/GameXost/Avito_Test_Case/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (u *UserRepo) GetPR(ctx context.Context, userID string) (string, []models.PullRequestShort) {
	var PRequests []models.PullRequestShort
	query := `
	SELECT pull_request.pull_request_id, pull_request.pull_request_name, pull_request.author_id, pull_request.status
	FROM users
	JOIN pull_request ON users.user_id = pull_request.author_id
	WHERE users.user_id = $1`

	rows, err := u.pool.Query(ctx, query, userID)
	if err != nil {
		return "", nil
	}
	defer rows.Close()
	for rows.Next() {
		var PRShort models.PullRequestShort
		err = rows.Scan(&PRShort.PullRequestID, &PRShort.PullRequestName, &PRShort.AuthorID, &PRShort.Status)
		if err != nil {
			return "", nil
		}
		PRequests = append(PRequests, PRShort)
	}
	err = rows.Err()
	if err != nil {
		return "", nil
	}
	return userID, PRequests
}
