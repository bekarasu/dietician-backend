package openai

type Meal struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Calories    int    `json:"calories" validate:"required"`
	Protein     int    `json:"protein" validate:"required"`
	Carbs       int    `json:"carbs" validate:"required"`
	Fats        int    `json:"fats" validate:"required"`
}

type DailyPlan struct {
	Day           int    `json:"day" validate:"required,min=1"`
	Breakfast     Meal   `json:"breakfast" validate:"required"`
	Lunch         Meal   `json:"lunch" validate:"required"`
	Dinner        Meal   `json:"dinner" validate:"required"`
	Snacks        []Meal `json:"snacks" validate:"dive"`
	TotalCalories int    `json:"totalCalories" validate:"required"`
}

type DietRecommendationResponse struct {
	DailyPlans []DailyPlan `json:"dailyPlans" validate:"required,min=1,dive"`
	Summary    string      `json:"summary" validate:"required"`
}
