package request

type AddWeightRequest struct {
	WeightKg float64 `json:"weightKg"`
	Notes    *string `json:"notes,omitempty"`
}

type AddHabitRequest struct {
	HabitName string  `json:"habitName"`
	Completed bool    `json:"completed"`
	Notes     *string `json:"notes,omitempty"`
}
