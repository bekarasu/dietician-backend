package model

import (
	"time"
)

type DietPlan struct {
	ID        string          `db:"id"`
	UserID    string          `db:"user_id"`
	StartDate time.Time       `db:"start_date"`
	EndDate   *time.Time      `db:"end_date"`
	Goals     *string         `db:"goals"`
	Status    string          `db:"status"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
	Meals     []DietPlanMeal `db:"-"`
}

type DietPlanMeal struct {
	ID          string     `db:"id"`
	DietPlanID  string     `db:"diet_plan_id"`
	DayOfWeek   *int       `db:"day_of_week"`
	MealType    string     `db:"meal_type"`
	RecipeID    *string    `db:"recipe_id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	Calories    *int       `db:"calories"`
	ProteinG    *float64   `db:"protein_g"`
	CarbsG      *float64   `db:"carbs_g"`
	FatG        *float64   `db:"fat_g"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
