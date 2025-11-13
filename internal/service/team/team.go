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
