package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/internal/repository"
	"github.com/GameXost/Avito_Test_Case/models"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) GetUserReviews(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	pullRequests, err := u.userRepo.GetUserReviews(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserReviews: %w", err)
	}
	return pullRequests, nil
}

func (u *UserService) SetActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	var user *models.User
	user, err := u.userRepo.SetActive(ctx, userID, isActive)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrNotFound
		default:
			return nil, fmt.Errorf("SetActive: %w", err)
		}
	}
	return user, nil
}
