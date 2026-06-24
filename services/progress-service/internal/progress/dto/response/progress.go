package response

import (
	"dietician.local/services/progress-service/internal/progress/repository"
	dailylog_dto "dietician.local/services/progress-service/internal/dailylog/dto"
	tracking_dto "dietician.local/services/progress-service/internal/tracking/dto"
)

type ProgressResponse struct {
	WeightLogs      []repository.WeightLog                `json:"weightLogs"`
	HabitLogs       []repository.HabitLog                 `json:"habitLogs"`
	DailyLogs       []dailylog_dto.DailyLogResponse       `json:"dailyLogs"`
	TrackingMetrics []tracking_dto.TrackingMetricResponse `json:"trackingMetrics"`
}
