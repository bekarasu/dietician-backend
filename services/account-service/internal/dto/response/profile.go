package response

import "time"

type UserProfile struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	DateOfBirth   *string   `json:"date_of_birth,omitempty"`
	Gender        *string   `json:"gender,omitempty"`
	HeightCm      *float64  `json:"height_cm,omitempty"`
	WeightKg      *float64  `json:"weight_kg,omitempty"`
	ActivityLevel *string   `json:"activity_level,omitempty"`
	Goal          *string   `json:"goal,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DietaryPreference struct {
	ID         string    `json:"id"`
	Preference string    `json:"preference"`
	CreatedAt  time.Time `json:"created_at"`
}

type Allergy struct {
	ID        string    `json:"id"`
	Allergy   string    `json:"allergy"`
	Severity  *string   `json:"severity,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type DislikedFood struct {
	ID        string    `json:"id"`
	FoodName  string    `json:"food_name"`
	CreatedAt time.Time `json:"created_at"`
}

type PreferencesResponse struct {
	Preferences   []DietaryPreference `json:"preferences"`
	Allergies     []Allergy           `json:"allergies"`
	DislikedFoods []DislikedFood      `json:"disliked_foods"`
}
