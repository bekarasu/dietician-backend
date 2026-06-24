package dto

import "time"

type CreateDietPlanRequest struct {
	UserID    string                  `json:"user_id" validate:"required"`
	StartDate time.Time               `json:"start_date" validate:"required"`
	EndDate   *time.Time              `json:"end_date"`
	Goals     *string                 `json:"goals"`
	Meals     []CreateDietPlanMealDto `json:"meals" validate:"required,min=1"`
}

type CreateDietPlanMealDto struct {
	DayOfWeek   *int     `json:"day_of_week"`
	MealType    string   `json:"meal_type" validate:"required"`
	RecipeID    *string  `json:"recipe_id"`
	Name        string   `json:"name" validate:"required"`
	Description *string  `json:"description"`
	Calories    *int     `json:"calories"`
	ProteinG    *float64 `json:"protein_g"`
	CarbsG      *float64 `json:"carbs_g"`
	FatG        *float64 `json:"fat_g"`
}

type DietPlanResponse struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	StartDate time.Time              `json:"start_date"`
	EndDate   *time.Time             `json:"end_date"`
	Goals     *string                `json:"goals"`
	Status    string                 `json:"status"`
	Meals     []DietPlanMealResponse `json:"meals"`
}

type DietPlanMealResponse struct {
	ID          string   `json:"id"`
	DayOfWeek   *int     `json:"day_of_week"`
	MealType    string   `json:"meal_type"`
	RecipeID    *string  `json:"recipe_id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Calories    *int     `json:"calories"`
	ProteinG    *float64 `json:"protein_g"`
	CarbsG      *float64 `json:"carbs_g"`
	FatG        *float64 `json:"fat_g"`
}
