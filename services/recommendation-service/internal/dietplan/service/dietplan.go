package service

import (
	"context"

	"dietician.local/services/recommendation-service/internal/dietplan/dto"
	"dietician.local/services/recommendation-service/internal/dietplan/model"
	"dietician.local/services/recommendation-service/internal/dietplan/repository"
)

type IDietPlanService interface {
	CreateDietPlan(ctx context.Context, req dto.CreateDietPlanRequest) (*dto.DietPlanResponse, error)
}

type dietPlanService struct {
	repo repository.IDietPlanRepository
}

func NewDietPlanService(repo repository.IDietPlanRepository) IDietPlanService {
	return &dietPlanService{
		repo: repo,
	}
}

func (s *dietPlanService) CreateDietPlan(ctx context.Context, req dto.CreateDietPlanRequest) (*dto.DietPlanResponse, error) {
	plan := &model.DietPlan{
		UserID:    req.UserID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Goals:     req.Goals,
		Status:    "active",
	}

	for _, mealReq := range req.Meals {
		meal := model.DietPlanMeal{
			DayOfWeek:   mealReq.DayOfWeek,
			MealType:    mealReq.MealType,
			RecipeID:    mealReq.RecipeID,
			Name:        mealReq.Name,
			Description: mealReq.Description,
			Calories:    mealReq.Calories,
			ProteinG:    mealReq.ProteinG,
			CarbsG:      mealReq.CarbsG,
			FatG:        mealReq.FatG,
		}
		plan.Meals = append(plan.Meals, meal)
	}

	err := s.repo.CreateDietPlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	res := &dto.DietPlanResponse{
		ID:        plan.ID,
		UserID:    plan.UserID,
		StartDate: plan.StartDate,
		EndDate:   plan.EndDate,
		Goals:     plan.Goals,
		Status:    plan.Status,
	}

	for _, meal := range plan.Meals {
		res.Meals = append(res.Meals, dto.DietPlanMealResponse{
			ID:          meal.ID,
			DayOfWeek:   meal.DayOfWeek,
			MealType:    meal.MealType,
			RecipeID:    meal.RecipeID,
			Name:        meal.Name,
			Description: meal.Description,
			Calories:    meal.Calories,
			ProteinG:    meal.ProteinG,
			CarbsG:      meal.CarbsG,
			FatG:        meal.FatG,
		})
	}

	return res, nil
}
