package service

import (
	"context"
	"time"

	"dietician.local/services/progress-service/internal/dailylog/dto"
	"dietician.local/services/progress-service/internal/dailylog/model"
	"dietician.local/services/progress-service/internal/dailylog/repository"
)

type IDailyLogService interface {
	CreateDailyLog(ctx context.Context, req dto.CreateDailyLogRequest) (*dto.DailyLogResponse, error)
	GetDailyLog(ctx context.Context, userID string, date string) (*dto.DailyLogResponse, error)
	GetDailyLogs(ctx context.Context, userID string) ([]dto.DailyLogResponse, error)
}

type dailyLogService struct {
	repo repository.IDailyLogRepository
}

func NewDailyLogService(repo repository.IDailyLogRepository) IDailyLogService {
	return &dailyLogService{
		repo: repo,
	}
}

func (s *dailyLogService) CreateDailyLog(ctx context.Context, req dto.CreateDailyLogRequest) (*dto.DailyLogResponse, error) {
	logDate, _ := time.Parse("2006-01-02", req.LogDate)
	log := &model.DailyLog{
		UserID:         req.UserID,
		LogDate:        logDate,
		WaterIntakeML:  req.WaterIntakeML,
		SleepHours:     req.SleepHours,
		ExerciseMins:   req.ExerciseMins,
		Mood:           req.Mood,
		Notes:          req.Notes,
	}

	for _, mealReq := range req.Meals {
		log.Meals = append(log.Meals, model.DailyLogMeal{
			MealType: mealReq.MealType,
			Name:     mealReq.Name,
			Calories: mealReq.Calories,
			ProteinG: mealReq.ProteinG,
			CarbsG:   mealReq.CarbsG,
			FatG:     mealReq.FatG,
		})
	}

	err := s.repo.CreateDailyLog(ctx, log)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(log), nil
}

func (s *dailyLogService) GetDailyLog(ctx context.Context, userID string, date string) (*dto.DailyLogResponse, error) {
	log, err := s.repo.GetDailyLog(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(log), nil
}

func (s *dailyLogService) GetDailyLogs(ctx context.Context, userID string) ([]dto.DailyLogResponse, error) {
	logs, err := s.repo.GetDailyLogs(ctx, userID)
	if err != nil {
		return nil, err
	}
	var res []dto.DailyLogResponse
	for _, l := range logs {
		res = append(res, *s.mapToResponse(&l))
	}
	return res, nil
}

func (s *dailyLogService) mapToResponse(log *model.DailyLog) *dto.DailyLogResponse {
	res := &dto.DailyLogResponse{
		ID:             log.ID,
		UserID:         log.UserID,
		LogDate:        log.LogDate,
		WaterIntakeML:  log.WaterIntakeML,
		SleepHours:     log.SleepHours,
		ExerciseMins:   log.ExerciseMins,
		Mood:           log.Mood,
		Notes:          log.Notes,
	}
	for _, m := range log.Meals {
		res.Meals = append(res.Meals, dto.DailyLogMealRes{
			ID:       m.ID,
			MealType: m.MealType,
			Name:     m.Name,
			Calories: m.Calories,
			ProteinG: m.ProteinG,
			CarbsG:   m.CarbsG,
			FatG:     m.FatG,
			LoggedAt: m.LoggedAt,
		})
	}
	return res
}
