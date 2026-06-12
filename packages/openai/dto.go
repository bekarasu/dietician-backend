package openai

type DietRecommendationRequestDto struct {
	Age                int      `json:"age" validate:"required,min=1"`
	Gender             string   `json:"gender" validate:"required"`
	HeightCm           float64  `json:"heightCm" validate:"required,min=50"`
	WeightKg           float64  `json:"weightKg" validate:"required,min=20"`
	ActivityLevel      string   `json:"activityLevel" validate:"required"` // e.g. sedentary, light, moderate, active, very active
	Goal               string   `json:"goal" validate:"required"`          // e.g. weight_loss, muscle_gain, maintenance
	DietaryPreferences []string `json:"dietaryPreferences"`                // e.g. vegan, keto
	Allergies          []string `json:"allergies"`
	DurationDays       int      `json:"durationDays" validate:"required,min=1,max=7"`
}
