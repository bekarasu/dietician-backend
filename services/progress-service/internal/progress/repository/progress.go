package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type WeightLog struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	WeightKg  float64   `db:"weight_kg" json:"weightKg"`
	Notes     *string   `db:"notes" json:"notes,omitempty"`
	LoggedAt  time.Time `db:"logged_at" json:"loggedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type HabitLog struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	HabitName string    `db:"habit_name" json:"habitName"`
	Completed bool      `db:"completed" json:"completed"`
	Notes     *string   `db:"notes" json:"notes,omitempty"`
	LoggedAt  time.Time `db:"logged_at" json:"loggedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type WeeklyProgressSummary struct {
	ID              string    `db:"id" json:"id"`
	UserID          string    `db:"user_id" json:"userId"`
	WeekStart       string    `db:"week_start" json:"weekStart"`
	WeekEnd         string    `db:"week_end" json:"weekEnd"`
	StartWeight     *float64  `db:"start_weight" json:"startWeight,omitempty"`
	EndWeight       *float64  `db:"end_weight" json:"endWeight,omitempty"`
	WeightChange    *float64  `db:"weight_change" json:"weightChange,omitempty"`
	HabitsCompleted int       `db:"habits_completed" json:"habitsCompleted"`
	HabitsTotal     int       `db:"habits_total" json:"habitsTotal"`
	Summary         *string   `db:"summary" json:"summary,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

type IProgressRepository interface {
	AddWeightLog(ctx context.Context, log *WeightLog) error
	GetWeightLogs(ctx context.Context, userID string) ([]WeightLog, error)
	AddHabitLog(ctx context.Context, log *HabitLog) error
	GetHabitLogs(ctx context.Context, userID string) ([]HabitLog, error)
	GetWeeklySummary(ctx context.Context, userID string) (*WeeklyProgressSummary, error)
}

type progressRepository struct {
	db *sqlx.DB
}

func NewProgressRepository(db *sqlx.DB) IProgressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) AddWeightLog(ctx context.Context, log *WeightLog) error {
	query := `INSERT INTO weight_logs (user_id, weight_kg, notes)
		VALUES ($1, $2, $3) RETURNING id, logged_at, created_at`
	return r.db.QueryRowContext(ctx, query,
		log.UserID, log.WeightKg, log.Notes,
	).Scan(&log.ID, &log.LoggedAt, &log.CreatedAt)
}

func (r *progressRepository) GetWeightLogs(ctx context.Context, userID string) ([]WeightLog, error) {
	var logs []WeightLog
	err := r.db.SelectContext(ctx, &logs,
		`SELECT id, user_id, weight_kg, notes, logged_at, created_at
		FROM weight_logs WHERE user_id = $1 ORDER BY logged_at DESC`, userID)
	return logs, err
}

func (r *progressRepository) AddHabitLog(ctx context.Context, log *HabitLog) error {
	query := `INSERT INTO habit_logs (user_id, habit_name, completed, notes)
		VALUES ($1, $2, $3, $4) RETURNING id, logged_at, created_at`
	return r.db.QueryRowContext(ctx, query,
		log.UserID, log.HabitName, log.Completed, log.Notes,
	).Scan(&log.ID, &log.LoggedAt, &log.CreatedAt)
}

func (r *progressRepository) GetHabitLogs(ctx context.Context, userID string) ([]HabitLog, error) {
	var logs []HabitLog
	err := r.db.SelectContext(ctx, &logs,
		`SELECT id, user_id, habit_name, completed, notes, logged_at, created_at
		FROM habit_logs WHERE user_id = $1 ORDER BY logged_at DESC`, userID)
	return logs, err
}

func (r *progressRepository) GetWeeklySummary(ctx context.Context, userID string) (*WeeklyProgressSummary, error) {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 6)

	// Get weight logs for this week
	var weights []WeightLog
	err := r.db.SelectContext(ctx, &weights,
		`SELECT id, user_id, weight_kg, notes, logged_at, created_at
		FROM weight_logs WHERE user_id = $1 AND logged_at >= $2 AND logged_at <= $3
		ORDER BY logged_at ASC`, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	// Get habit logs for this week
	var habits []HabitLog
	err = r.db.SelectContext(ctx, &habits,
		`SELECT id, user_id, habit_name, completed, notes, logged_at, created_at
		FROM habit_logs WHERE user_id = $1 AND logged_at >= $2 AND logged_at <= $3`,
		userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	summary := &WeeklyProgressSummary{
		UserID:      userID,
		WeekStart:   weekStart.Format("2006-01-02"),
		WeekEnd:     weekEnd.Format("2006-01-02"),
		HabitsTotal: len(habits),
	}

	if len(weights) > 0 {
		sw := weights[0].WeightKg
		ew := weights[len(weights)-1].WeightKg
		wc := ew - sw
		summary.StartWeight = &sw
		summary.EndWeight = &ew
		summary.WeightChange = &wc
	}

	for _, h := range habits {
		if h.Completed {
			summary.HabitsCompleted++
		}
	}

	return summary, nil
}
