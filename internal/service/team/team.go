package service

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/request"
)

func TeamFromRequest(teamRequest request.TeamRequest) *domain.Team {
	return &domain.Team{
		TeamName: teamRequest.TeamName,
	}
}

func TeamMembersFromRequest(teamId string, teamMembers []request.TeamMember) []domain.TeamMember {
	var members []domain.TeamMember

	for _, teamMember := range teamMembers {
		member := domain.TeamMember{
			TeamID:   teamId,
			UserID:   teamMember.UserId,
			Username: teamMember.Username,
			IsActive: teamMember.IsActive,
		}
		members = append(members, member)
	}

	return members
}
