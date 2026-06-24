package model

import "time"

type DailyLog struct {
	ID             string          `db:"id"`
	UserID         string          `db:"user_id"`
	LogDate        time.Time       `db:"log_date"`
	WaterIntakeML  int             `db:"water_intake_ml"`
	SleepHours     float64         `db:"sleep_hours"`
	ExerciseMins   int             `db:"exercise_minutes"`
	Mood           *string         `db:"mood"`
	Notes          *string         `db:"notes"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
	Meals          []DailyLogMeal `db:"-"`
}

type DailyLogMeal struct {
	ID         string    `db:"id"`
	DailyLogID string    `db:"daily_log_id"`
	MealType   string    `db:"meal_type"`
	Name       string    `db:"name"`
	Calories   *int      `db:"calories"`
	ProteinG   *float64  `db:"protein_g"`
	CarbsG     *float64  `db:"carbs_g"`
	FatG       *float64  `db:"fat_g"`
	LoggedAt   time.Time `db:"logged_at"`
}
