package dto

import "time"

type CreateDailyLogRequest struct {
	UserID         string               `json:"user_id" validate:"required"`
	LogDate        string               `json:"log_date" validate:"required"` // Format: YYYY-MM-DD
	WaterIntakeML  int                  `json:"water_intake_ml"`
	SleepHours     float64              `json:"sleep_hours"`
	ExerciseMins   int                  `json:"exercise_minutes"`
	Mood           *string              `json:"mood"`
	Notes          *string              `json:"notes"`
	Meals          []CreateDailyLogMeal `json:"meals"`
}

type CreateDailyLogMeal struct {
	MealType   string   `json:"meal_type" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	Calories   *int     `json:"calories"`
	ProteinG   *float64 `json:"protein_g"`
	CarbsG     *float64 `json:"carbs_g"`
	FatG       *float64 `json:"fat_g"`
}

type DailyLogResponse struct {
	ID             string               `json:"id"`
	UserID         string               `json:"user_id"`
	LogDate        time.Time            `json:"log_date"`
	WaterIntakeML  int                  `json:"water_intake_ml"`
	SleepHours     float64              `json:"sleep_hours"`
	ExerciseMins   int                  `json:"exercise_minutes"`
	Mood           *string              `json:"mood"`
	Notes          *string              `json:"notes"`
	Meals          []DailyLogMealRes    `json:"meals"`
}

type DailyLogMealRes struct {
	ID         string    `json:"id"`
	MealType   string    `json:"meal_type"`
	Name       string    `json:"name"`
	Calories   *int      `json:"calories"`
	ProteinG   *float64  `json:"protein_g"`
	CarbsG     *float64  `json:"carbs_g"`
	FatG       *float64  `json:"fat_g"`
	LoggedAt   time.Time `json:"logged_at"`
}
