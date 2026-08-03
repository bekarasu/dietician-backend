package repository

import (
	"context"
	"database/sql"

	"dietician.local/services/account-service/internal/profile/model"
	"github.com/jmoiron/sqlx"
)

type ProfileRepository struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID string) (*model.UserProfile, error) {
	var profile model.UserProfile
	query := `SELECT id, user_id, date_of_birth, gender, height_cm, weight_kg, activity_level, goal, 
		age, display_name, target_weight_kg, daily_calorie_target, target_water_ml, target_coffee_cups,
		created_at, updated_at
		FROM user_profiles WHERE user_id = $1`
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &profile, err
}

func (r *ProfileRepository) Upsert(ctx context.Context, profile *model.UserProfile) error {
	query := `INSERT INTO user_profiles (user_id, date_of_birth, gender, height_cm, weight_kg, activity_level, goal, age, display_name, target_weight_kg, daily_calorie_target, target_water_ml, target_coffee_cups)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id) DO UPDATE SET
			date_of_birth = COALESCE(EXCLUDED.date_of_birth, user_profiles.date_of_birth),
			gender = COALESCE(EXCLUDED.gender, user_profiles.gender),
			height_cm = COALESCE(EXCLUDED.height_cm, user_profiles.height_cm),
			weight_kg = COALESCE(EXCLUDED.weight_kg, user_profiles.weight_kg),
			activity_level = COALESCE(EXCLUDED.activity_level, user_profiles.activity_level),
			goal = COALESCE(EXCLUDED.goal, user_profiles.goal),
			age = COALESCE(EXCLUDED.age, user_profiles.age),
			display_name = COALESCE(EXCLUDED.display_name, user_profiles.display_name),
			target_weight_kg = COALESCE(EXCLUDED.target_weight_kg, user_profiles.target_weight_kg),
			daily_calorie_target = COALESCE(EXCLUDED.daily_calorie_target, user_profiles.daily_calorie_target),
			target_water_ml = COALESCE(EXCLUDED.target_water_ml, user_profiles.target_water_ml),
			target_coffee_cups = COALESCE(EXCLUDED.target_coffee_cups, user_profiles.target_coffee_cups),
			updated_at = NOW()
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		profile.UserID, profile.DateOfBirth, profile.Gender,
		profile.HeightCm, profile.WeightKg, profile.ActivityLevel, profile.Goal,
		profile.Age, profile.DisplayName, profile.TargetWeightKg, profile.DailyCalorieTarget,
		profile.TargetWaterMl, profile.TargetCoffeeCups,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
}

func (r *ProfileRepository) GetPreferences(ctx context.Context, userID string) ([]model.DietaryPreference, error) {
	var prefs []model.DietaryPreference
	err := r.db.SelectContext(ctx, &prefs,
		`SELECT id, user_id, preference, created_at FROM dietary_preferences WHERE user_id = $1`, userID)
	return prefs, err
}

func (r *ProfileRepository) ReplacePreferences(ctx context.Context, userID string, preferences []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM dietary_preferences WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, p := range preferences {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO dietary_preferences (user_id, preference) VALUES ($1, $2)`, userID, p); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Allergies

func (r *ProfileRepository) GetAllergies(ctx context.Context, userID string) ([]model.Allergy, error) {
	var allergies []model.Allergy
	err := r.db.SelectContext(ctx, &allergies,
		`SELECT id, user_id, allergy, severity, created_at FROM allergies WHERE user_id = $1`, userID)
	return allergies, err
}

func (r *ProfileRepository) ReplaceAllergies(ctx context.Context, userID string, allergies []model.AllergyInput) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM allergies WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, a := range allergies {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO allergies (user_id, allergy, severity) VALUES ($1, $2, $3)`,
			userID, a.Allergy, a.Severity); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Disliked Foods

func (r *ProfileRepository) GetDislikedFoods(ctx context.Context, userID string) ([]model.DislikedFood, error) {
	var foods []model.DislikedFood
	err := r.db.SelectContext(ctx, &foods,
		`SELECT id, user_id, food_name, created_at FROM disliked_foods WHERE user_id = $1`, userID)
	return foods, err
}

func (r *ProfileRepository) ReplaceDislikedFoods(ctx context.Context, userID string, foods []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM disliked_foods WHERE user_id = $1`, userID); err != nil {
		return err
	}
	for _, f := range foods {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO disliked_foods (user_id, food_name) VALUES ($1, $2)`, userID, f); err != nil {
			return err
		}
	}
	return tx.Commit()
}
