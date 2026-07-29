package dto

import "time"

type CreateDailyLogRequest struct {
	UserID         string               `json:"userId" validate:"required"`
	LogDate        string               `json:"logDate" validate:"required"` // Format: YYYY-MM-DD
	WaterIntakeML  int                  `json:"waterIntakeMl"`
	SleepHours     float64              `json:"sleepHours"`
	ExerciseMins   int                  `json:"exerciseMinutes"`
	Mood           *string              `json:"mood"`
	Notes          *string              `json:"notes"`
	Meals          []CreateDailyLogMeal `json:"meals"`
}

type CreateDailyLogMeal struct {
	MealType   string   `json:"mealType" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	Calories   *int     `json:"calories"`
	ProteinG   *float64 `json:"proteinG"`
	CarbsG     *float64 `json:"carbsG"`
	FatG       *float64 `json:"fatG"`
}

type DailyLogResponse struct {
	ID             string               `json:"id"`
	UserID         string               `json:"userId"`
	LogDate        time.Time            `json:"logDate"`
	WaterIntakeML  int                  `json:"waterIntakeMl"`
	SleepHours     float64              `json:"sleepHours"`
	ExerciseMins   int                  `json:"exerciseMinutes"`
	Mood           *string              `json:"mood"`
	Notes          *string              `json:"notes"`
	Meals          []DailyLogMealRes    `json:"meals"`
}

type DailyLogMealRes struct {
	ID         string    `json:"id"`
	MealType   string    `json:"mealType"`
	Name       string    `json:"name"`
	Calories   *int      `json:"calories"`
	ProteinG   *float64  `json:"proteinG"`
	CarbsG     *float64  `json:"carbsG"`
	FatG       *float64  `json:"fatG"`
	LoggedAt   time.Time `json:"loggedAt"`
}
