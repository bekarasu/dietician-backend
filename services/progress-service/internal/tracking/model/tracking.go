package model

import "time"

type TrackingMetric struct {
	ID           string     `db:"id"`
	UserID       string     `db:"user_id"`
	MetricType   string     `db:"metric_type"`
	TargetValue  *float64   `db:"target_value"`
	CurrentValue float64    `db:"current_value"`
	Deadline     *time.Time `db:"deadline"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
