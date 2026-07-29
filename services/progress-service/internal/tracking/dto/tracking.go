package dto

import "time"

type CreateTrackingMetricRequest struct {
	UserID       string     `json:"userId" validate:"required"`
	MetricType   string     `json:"metricType" validate:"required"`
	TargetValue  *float64   `json:"targetValue"`
	CurrentValue float64    `json:"currentValue" validate:"required"`
	Deadline     *time.Time `json:"deadline"`
}

type TrackingMetricResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	MetricType   string     `json:"metricType"`
	TargetValue  *float64   `json:"targetValue"`
	CurrentValue float64    `json:"currentValue"`
	Deadline     *time.Time `json:"deadline"`
	Status       string     `json:"status"`
}
