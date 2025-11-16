package teamService

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/request"
	"context"
	"fmt"
)

type TeamAddStorage interface {
	SaveTeamWithMembers(ctx context.Context, team *domain.Team, members []domain.TeamMember) (string, error)
}

type TeamAddService struct {
	storage TeamAddStorage
}

func NewTeamAddService(storage TeamAddStorage) *TeamAddService {
	return &TeamAddService{storage: storage}
}

func (s *TeamAddService) AddTeam(ctx context.Context, req *request.TeamAddRequest) (string, error) {
	const op = "service.TeamAddService.AddTeam"

	team := TeamFromRequest(*req)
	members := TeamMembersFromRequest(req.Members)

	teamId, err := s.storage.SaveTeamWithMembers(ctx, team, members)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return teamId, nil
}

func TeamFromRequest(teamRequest request.TeamAddRequest) *domain.Team {
	return &domain.Team{
		TeamName: teamRequest.TeamName,
	}
}

func TeamMembersFromRequest(teamMembers []request.TeamMember) []domain.TeamMember {
	var members []domain.TeamMember

	for _, teamMember := range teamMembers {
		member := domain.TeamMember{
			UserID:   teamMember.UserId,
			Username: teamMember.Username,
			IsActive: teamMember.IsActive,
		}
		members = append(members, member)
	}

	return members
}
