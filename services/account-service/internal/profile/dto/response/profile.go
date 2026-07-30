package response

import "time"

type UserProfile struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"userId"`
	DateOfBirth        *string   `json:"dateOfBirth,omitempty"`
	Gender             *string   `json:"gender,omitempty"`
	HeightCm           *float64  `json:"heightCm,omitempty"`
	WeightKg           *float64  `json:"weightKg,omitempty"`
	ActivityLevel      *string   `json:"activityLevel,omitempty"`
	Goal               *string   `json:"goal,omitempty"`
	Age                *int      `json:"age,omitempty"`
	DisplayName        *string   `json:"displayName,omitempty"`
	TargetWeightKg     *float64  `json:"targetWeightKg,omitempty"`
	DailyCalorieTarget *int      `json:"dailyCalorieTarget,omitempty"`
	TargetWaterMl      *int      `json:"targetWaterMl,omitempty"`
	TargetCoffeeCups   *int      `json:"targetCoffeeCups,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type DietaryPreference struct {
	ID         string    `json:"id"`
	Preference string    `json:"preference"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Allergy struct {
	ID        string    `json:"id"`
	Allergy   string    `json:"allergy"`
	Severity  *string   `json:"severity,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type DislikedFood struct {
	ID        string    `json:"id"`
	FoodName  string    `json:"foodName"`
	CreatedAt time.Time `json:"createdAt"`
}

type PreferencesResponse struct {
	Preferences   []DietaryPreference `json:"preferences"`
	Allergies     []Allergy           `json:"allergies"`
	DislikedFoods []DislikedFood      `json:"dislikedFoods"`
}
