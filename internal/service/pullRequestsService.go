package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/Avito_Test_Case/internal/repository"
	"github.com/GameXost/Avito_Test_Case/models"
	"slices"
)

type PullRequestService struct {
	prRepo   *repository.PullRequestRepo
	teamRepo *repository.TeamRepo
}

func NewPRService(prRepo *repository.PullRequestRepo, teamRepo *repository.TeamRepo) *PullRequestService {
	return &PullRequestService{
		prRepo:   prRepo,
		teamRepo: teamRepo,
	}
}

func (p *PullRequestService) CreatePR(ctx context.Context, prID, prName, authorID string) (*models.PullRequest, error) {
	pullRequest := models.PullRequest{PullRequestID: prID, PullRequestName: prName, AuthorID: authorID, Status: "OPEN"}

	teamName, err := p.prRepo.GetTeamNameByUserID(ctx, authorID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrNotFound
		default:
			return nil, fmt.Errorf("CreatePR: %w", err)
		}
	}

	team, err := p.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("CreatePR: %w", err)
	}

	for _, member := range team.Members {
		if len(pullRequest.AssignedReviewers) == 2 {
			break
		}

		if member.IsActive == true && member.UserId != authorID {
			pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, member.UserId)
		}
	}

	err = p.prRepo.CreatePR(ctx, pullRequest)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrNotFound
		case errors.Is(err, models.ErrPRExists):
			return nil, models.ErrPRExists
		default:
			return nil, fmt.Errorf("CreatePR: %w", err)
		}
	}
	return &pullRequest, nil
}

func (p *PullRequestService) Merge(ctx context.Context, prID string) (*models.PullRequest, error) {
	pullRequest, err := p.prRepo.MergePR(ctx, prID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("Merge: %w", err)
	}
	return pullRequest, nil
}

func (p *PullRequestService) Reassign(ctx context.Context, prID, oldRevID string) (*models.PullRequest, error) {
	status, err := p.prRepo.IsMerged(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("Reassign: %w", err)
	}
	if status == "MERGED" {
		return nil, models.ErrPRMerged
	}

	// Получаем PR и его ревьюверов
	pullRequest, err := p.prRepo.GetPRInfo(ctx, prID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrNotFound
		default:
			return nil, fmt.Errorf("Reassign: %w", err)
		}
	}

	prReviewers, err := p.prRepo.GetReviewers(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("Reassign: %w", err)
	}
	pullRequest.AssignedReviewers = prReviewers

	if !slices.Contains(pullRequest.AssignedReviewers, oldRevID) {
		return nil, models.ErrNotAssigned
	}

	teamName, err := p.prRepo.GetTeamNameByUserID(ctx, oldRevID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrNotFound
		default:
			return nil, fmt.Errorf("Reassign: %w", err)
		}
	}

	team, err := p.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("Reassign: %w", err)
	}

	var newRev string
	for _, member := range team.Members {
		// Проверяем: активен ли член команды
		if !member.IsActive {
			continue
		}

		if !slices.Contains(pullRequest.AssignedReviewers, member.UserId) &&
			member.UserId != pullRequest.AuthorID {
			newRev = member.UserId

			err = p.prRepo.DelAndAssign(ctx, oldRevID, prID, member.UserId)
			if err != nil {
				switch {
				case errors.Is(err, models.ErrNotAssigned):
					return nil, models.ErrNotAssigned
				case errors.Is(err, models.ErrNotFound):
					return nil, models.ErrNotFound
				default:
					return nil, fmt.Errorf("Reassign: %w", err)
				}
			}

			pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, member.UserId)
			idx := slices.Index(pullRequest.AssignedReviewers, oldRevID)
			if idx != -1 {
				pullRequest.AssignedReviewers = slices.Delete(pullRequest.AssignedReviewers, idx, idx+1)
			}
			break
		}
	}

	if newRev == "" {
		return nil, models.ErrNoCandidate
	}

	return pullRequest, nil
}
