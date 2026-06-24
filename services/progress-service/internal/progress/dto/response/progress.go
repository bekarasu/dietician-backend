package response

import "dietician.local/services/progress-service/internal/progress/repository"

type ProgressResponse struct {
	WeightLogs []repository.WeightLog `json:"weightLogs"`
	HabitLogs  []repository.HabitLog  `json:"habitLogs"`
}
