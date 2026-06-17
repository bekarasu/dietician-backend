package service

import (
	"context"
	"errors"

	"dietician.local/services/progress-service/internal/progress/repository"
)

type IProgressService interface {
	GetProgress(ctx context.Context, userID string) (*ProgressResponse, error)
	AddWeight(ctx context.Context, userID string, req AddWeightRequest) (*repository.WeightLog, error)
	GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error)
	AddHabit(ctx context.Context, userID string, req AddHabitRequest) (*repository.HabitLog, error)
}

type AddWeightRequest struct {
	WeightKg float64 `json:"weightKg"`
	Notes    *string `json:"notes,omitempty"`
}

type AddHabitRequest struct {
	HabitName string  `json:"habitName"`
	Completed bool    `json:"completed"`
	Notes     *string `json:"notes,omitempty"`
}

type ProgressResponse struct {
	WeightLogs []repository.WeightLog `json:"weightLogs"`
	HabitLogs  []repository.HabitLog  `json:"habitLogs"`
}

type progressService struct {
	repo repository.IProgressRepository
}

func NewProgressService(repo repository.IProgressRepository) IProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) GetProgress(ctx context.Context, userID string) (*ProgressResponse, error) {
	weights, err := s.repo.GetWeightLogs(ctx, userID)
	if err != nil {
		return nil, err
	}
	habits, err := s.repo.GetHabitLogs(ctx, userID)
	if err != nil {
		return nil, err
	}

	if weights == nil {
		weights = []repository.WeightLog{}
	}
	if habits == nil {
		habits = []repository.HabitLog{}
	}

	return &ProgressResponse{
		WeightLogs: weights,
		HabitLogs:  habits,
	}, nil
}

func (s *progressService) AddWeight(ctx context.Context, userID string, req AddWeightRequest) (*repository.WeightLog, error) {
	if req.WeightKg <= 0 {
		return nil, errors.New("weight must be positive")
	}

	log := &repository.WeightLog{
		UserID:   userID,
		WeightKg: req.WeightKg,
		Notes:    req.Notes,
	}
	if err := s.repo.AddWeightLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *progressService) GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error) {
	return s.repo.GetWeeklySummary(ctx, userID)
}

func (s *progressService) AddHabit(ctx context.Context, userID string, req AddHabitRequest) (*repository.HabitLog, error) {
	if req.HabitName == "" {
		return nil, errors.New("habit name is required")
	}

	log := &repository.HabitLog{
		UserID:    userID,
		HabitName: req.HabitName,
		Completed: req.Completed,
		Notes:     req.Notes,
	}
	if err := s.repo.AddHabitLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}
