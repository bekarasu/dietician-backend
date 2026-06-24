package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"dietician.local/services/recommendation-service/internal/dietplan/model"
)

type IDietPlanRepository interface {
	CreateDietPlan(ctx context.Context, plan *model.DietPlan) error
	GetDietPlanByID(ctx context.Context, id string) (*model.DietPlan, error)
}

type dietPlanRepository struct {
	db *sqlx.DB
}

func NewDietPlanRepository(db *sqlx.DB) IDietPlanRepository {
	return &dietPlanRepository{
		db: db,
	}
}

func (r *dietPlanRepository) CreateDietPlan(ctx context.Context, plan *model.DietPlan) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO diet_plans (user_id, start_date, end_date, goals, status) 
			  VALUES (:user_id, :start_date, :end_date, :goals, :status) RETURNING id, created_at, updated_at`
	
	rows, err := tx.NamedQuery(query, plan)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
	}
	rows.Close()

	if len(plan.Meals) > 0 {
		for i := range plan.Meals {
			plan.Meals[i].DietPlanID = plan.ID
		}
		mealQuery := `INSERT INTO diet_plan_meals (diet_plan_id, day_of_week, meal_type, recipe_id, name, description, calories, protein_g, carbs_g, fat_g)
					  VALUES (:diet_plan_id, :day_of_week, :meal_type, :recipe_id, :name, :description, :calories, :protein_g, :carbs_g, :fat_g) RETURNING id`
		mealRows, err := tx.NamedQuery(mealQuery, plan.Meals)
		if err != nil {
			return err
		}
		
		i := 0
		for mealRows.Next() {
			mealRows.Scan(&plan.Meals[i].ID)
			i++
		}
		mealRows.Close()
	}

	return tx.Commit()
}

func (r *dietPlanRepository) GetDietPlanByID(ctx context.Context, id string) (*model.DietPlan, error) {
	var plan model.DietPlan
	err := r.db.GetContext(ctx, &plan, "SELECT * FROM diet_plans WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("diet plan not found")
		}
		return nil, err
	}

	err = r.db.SelectContext(ctx, &plan.Meals, "SELECT * FROM diet_plan_meals WHERE diet_plan_id = $1", id)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}
