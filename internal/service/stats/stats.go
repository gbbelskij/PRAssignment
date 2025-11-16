package statsService

import (
	"PRAssignment/internal/domain"
	"context"
)

type StatsService struct {
	storage StatsStorage
}

type StatsStorage interface {
	GetStats(ctx context.Context) (*domain.Stats, error)
}

func NewStatsService(storage StatsStorage) *StatsService {
	return &StatsService{storage: storage}
}

func (s *StatsService) GetStats(ctx context.Context) (*domain.Stats, error) {
	return s.storage.GetStats(ctx)
}
