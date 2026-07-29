package dto

import "time"

type CreateDietPlanRequest struct {
	UserID    string                  `json:"userId" validate:"required"`
	StartDate time.Time               `json:"startDate" validate:"required"`
	EndDate   *time.Time              `json:"endDate"`
	Goals     *string                 `json:"goals"`
	Meals     []CreateDietPlanMealDto `json:"meals" validate:"required,min=1"`
}

type CreateDietPlanMealDto struct {
	DayOfWeek   *int     `json:"dayOfWeek"`
	MealType    string   `json:"mealType" validate:"required"`
	RecipeID    *string  `json:"recipeId"`
	Name        string   `json:"name" validate:"required"`
	Description *string  `json:"description"`
	Calories    *int     `json:"calories"`
	ProteinG    *float64 `json:"proteinG"`
	CarbsG      *float64 `json:"carbsG"`
	FatG        *float64 `json:"fatG"`
}

type DietPlanResponse struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"userId"`
	StartDate time.Time              `json:"startDate"`
	EndDate   *time.Time             `json:"endDate"`
	Goals     *string                `json:"goals"`
	Status    string                 `json:"status"`
	Meals     []DietPlanMealResponse `json:"meals"`
}

type DietPlanMealResponse struct {
	ID          string   `json:"id"`
	DayOfWeek   *int     `json:"dayOfWeek"`
	MealType    string   `json:"mealType"`
	RecipeID    *string  `json:"recipeId"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Calories    *int     `json:"calories"`
	ProteinG    *float64 `json:"proteinG"`
	CarbsG      *float64 `json:"carbsG"`
	FatG        *float64 `json:"fatG"`
}
