package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"dietician.local/services/progress-service/internal/tracking/model"
)

type ITrackingRepository interface {
	CreateTrackingMetric(ctx context.Context, metric *model.TrackingMetric) error
	GetTrackingMetricsByUser(ctx context.Context, userID string) ([]model.TrackingMetric, error)
}

type trackingRepository struct {
	db *sqlx.DB
}

func NewTrackingRepository(db *sqlx.DB) ITrackingRepository {
	return &trackingRepository{db: db}
}

func (r *trackingRepository) CreateTrackingMetric(ctx context.Context, metric *model.TrackingMetric) error {
	query := `INSERT INTO tracking_metrics (user_id, metric_type, target_value, current_value, deadline, status) 
			  VALUES (:user_id, :metric_type, :target_value, :current_value, :deadline, :status)
			  RETURNING id, created_at, updated_at`
	
	rows, err := r.db.NamedQueryContext(ctx, query, metric)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&metric.ID, &metric.CreatedAt, &metric.UpdatedAt)
	}
	return nil
}

func (r *trackingRepository) GetTrackingMetricsByUser(ctx context.Context, userID string) ([]model.TrackingMetric, error) {
	var metrics []model.TrackingMetric
	err := r.db.SelectContext(ctx, &metrics, "SELECT * FROM tracking_metrics WHERE user_id = $1", userID)
	return metrics, err
}
