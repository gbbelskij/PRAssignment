package teamService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/response"
	"context"
	"fmt"
)

type TeamGetStorage interface {
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	GetMembers(ctx context.Context, teamID string) ([]domain.TeamMember, error)
}

type TeamGetService struct {
	storage TeamGetStorage
}

func NewTeamGetService(storage TeamGetStorage) *TeamGetService {
	return &TeamGetService{storage: storage}
}

func (s *TeamGetService) GetTeamWithMembers(ctx context.Context, teamName string) (*response.TeamGetResponse, error) {
	team, err := s.storage.GetTeam(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	members, err := s.storage.GetMembers(ctx, team.TeamID)
	if err != nil {
		return nil, fmt.Errorf("get members: %w", err)
	}

	responseMembers := make([]response.TeamMember, 0, len(members))
	for _, m := range members {
		responseMembers = append(responseMembers, response.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return &response.TeamGetResponse{
		TeamName: team.TeamName,
		Members:  responseMembers,
	}, nil
}
