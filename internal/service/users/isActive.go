package userService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type SetIsActiveStorage interface {
	UpdateUser(ctx context.Context, userId string, isActive bool) (*domain.TeamMember, error)
	GetTeamNameById(ctx context.Context, teamId string) (string, error)
}

type SetIsActiveService struct {
	storage SetIsActiveStorage
}

func NewSetIsActiveService(storage SetIsActiveStorage) *SetIsActiveService {
	return &SetIsActiveService{storage: storage}
}

func (s *SetIsActiveService) SetUserActiveStatus(ctx context.Context, userId string, isActive bool) (*response.UserSetIsActiveResponse, error) {
	const op = "service.SetIsActiveService.SetUserActiveStatus"

	teamMember, err := s.storage.UpdateUser(ctx, userId, isActive)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	teamName, err := s.storage.GetTeamNameById(ctx, teamMember.TeamID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp := response.MakeUserSetIsActiveResponse(
		teamMember.UserID,
		teamMember.Username,
		teamName,
		teamMember.IsActive,
	)

	return &resp, nil
}
