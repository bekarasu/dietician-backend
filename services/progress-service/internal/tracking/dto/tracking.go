package dto

import "time"

type CreateTrackingMetricRequest struct {
	UserID       string     `json:"user_id" validate:"required"`
	MetricType   string     `json:"metric_type" validate:"required"`
	TargetValue  *float64   `json:"target_value"`
	CurrentValue float64    `json:"current_value" validate:"required"`
	Deadline     *time.Time `json:"deadline"`
}

type TrackingMetricResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	MetricType   string     `json:"metric_type"`
	TargetValue  *float64   `json:"target_value"`
	CurrentValue float64    `json:"current_value"`
	Deadline     *time.Time `json:"deadline"`
	Status       string     `json:"status"`
}
