package repository

import (
	"context"
	"database/sql"
	"errors"

	"dietician.local/packages/constants"
	"dietician.local/packages/utils"
	"dietician.local/services/progress-service/internal/dailylog/model"
	"github.com/jmoiron/sqlx"
)

type IDailyLogRepository interface {
	CreateDailyLog(ctx context.Context, log *model.DailyLog) error
	GetDailyLog(ctx context.Context, userID, date string) (*model.DailyLog, error)
	GetDailyLogs(ctx context.Context, userID string) ([]model.DailyLog, error)
}

type dailyLogRepository struct {
	db *sqlx.DB
}

func NewDailyLogRepository(db *sqlx.DB) IDailyLogRepository {
	return &dailyLogRepository{
		db: db,
	}
}

func (r *dailyLogRepository) CreateDailyLog(ctx context.Context, log *model.DailyLog) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO daily_logs (user_id, log_date, water_intake_ml, sleep_hours, exercise_minutes, mood, notes) 
			  VALUES (:user_id, :log_date, :water_intake_ml, :sleep_hours, :exercise_minutes, :mood, :notes) 
			  ON CONFLICT (user_id, log_date) DO UPDATE SET 
			    water_intake_ml = EXCLUDED.water_intake_ml,
			    sleep_hours = EXCLUDED.sleep_hours,
			    exercise_minutes = EXCLUDED.exercise_minutes,
			    mood = EXCLUDED.mood,
			    notes = EXCLUDED.notes,
			    updated_at = NOW()
			  RETURNING id, created_at, updated_at`
	
	rows, err := tx.NamedQuery(query, log)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&log.ID, &log.CreatedAt, &log.UpdatedAt)
	}
	rows.Close()

	if len(log.Meals) > 0 {
		for i := range log.Meals {
			log.Meals[i].DailyLogID = log.ID
		}
		mealQuery := `INSERT INTO daily_log_meals (daily_log_id, meal_type, name, calories, protein_g, carbs_g, fat_g)
					  VALUES (:daily_log_id, :meal_type, :name, :calories, :protein_g, :carbs_g, :fat_g) RETURNING id`
		mealRows, err := tx.NamedQuery(mealQuery, log.Meals)
		if err != nil {
			return err
		}
		
		i := 0
		for mealRows.Next() {
			mealRows.Scan(&log.Meals[i].ID)
			i++
		}
		mealRows.Close()
	}

	return tx.Commit()
}

func (r *dailyLogRepository) GetDailyLog(ctx context.Context, userID, date string) (*model.DailyLog, error) {
	var log model.DailyLog
	err := r.db.GetContext(ctx, &log, "SELECT * FROM daily_logs WHERE user_id = $1 AND log_date = $2", userID, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.DailyLogNotFound))
		}
		return nil, err
	}

	err = r.db.SelectContext(ctx, &log.Meals, "SELECT * FROM daily_log_meals WHERE daily_log_id = $1", log.ID)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *dailyLogRepository) GetDailyLogs(ctx context.Context, userID string) ([]model.DailyLog, error) {
	var logs []model.DailyLog
	err := r.db.SelectContext(ctx, &logs, "SELECT * FROM daily_logs WHERE user_id = $1 ORDER BY log_date DESC LIMIT 30", userID)
	if err != nil {
		return nil, err
	}
	
	for i := range logs {
		r.db.SelectContext(ctx, &logs[i].Meals, "SELECT * FROM daily_log_meals WHERE daily_log_id = $1", logs[i].ID)
	}

	return logs, nil
}
