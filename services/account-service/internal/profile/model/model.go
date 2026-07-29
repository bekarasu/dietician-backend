package model

import (
	"time"

	"dietician.local/services/account-service/internal/profile/dto/response"
)

type UserProfile struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	DateOfBirth   *string   `db:"date_of_birth" json:"date_of_birth,omitempty"`
	Gender        *string   `db:"gender" json:"gender,omitempty"`
	HeightCm      *float64  `db:"height_cm" json:"height_cm,omitempty"`
	WeightKg      *float64  `db:"weight_kg" json:"weight_kg,omitempty"`
	ActivityLevel *string   `db:"activity_level" json:"activity_level,omitempty"`
	Goal               *string   `db:"goal" json:"goal,omitempty"`
	Age                *int      `db:"age" json:"age,omitempty"`
	DisplayName        *string   `db:"display_name" json:"display_name,omitempty"`
	TargetWeightKg     *float64  `db:"target_weight_kg" json:"target_weight_kg,omitempty"`
	DailyCalorieTarget *int      `db:"daily_calorie_target" json:"daily_calorie_target,omitempty"`
	TargetWaterMl      *int      `db:"target_water_ml" json:"target_water_ml,omitempty"`
	TargetCoffeeCups   *int      `db:"target_coffee_cups" json:"target_coffee_cups,omitempty"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type DietaryPreference struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	Preference string    `db:"preference" json:"preference"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Allergy struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Allergy   string    `db:"allergy" json:"allergy"`
	Severity  *string   `db:"severity" json:"severity,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DislikedFood struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	FoodName  string    `db:"food_name" json:"food_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UpdateProfileRequest struct {
	DateOfBirth   *string  `json:"date_of_birth,omitempty"`
	Gender        *string  `json:"gender,omitempty"`
	HeightCm      *float64 `json:"height_cm,omitempty"`
	WeightKg      *float64 `json:"weight_kg,omitempty"`
	ActivityLevel *string  `json:"activity_level,omitempty"`
	Goal               *string  `json:"goal,omitempty"`
	Age                *int     `json:"age,omitempty"`
	DisplayName        *string  `json:"display_name,omitempty"`
	TargetWeightKg     *float64 `json:"target_weight_kg,omitempty"`
	DailyCalorieTarget *int     `json:"daily_calorie_target,omitempty"`
	TargetWaterMl      *int     `json:"target_water_ml,omitempty"`
	TargetCoffeeCups   *int     `json:"target_coffee_cups,omitempty"`
}

type OnboardingRequest struct {
	Name               string   `json:"name"`
	Age                int      `json:"age"`
	HeightCm           float64  `json:"heightCm"`
	WeightKg           float64  `json:"weightKg"`
	TargetWeightKg     float64  `json:"targetWeightKg"`
	Gender             string   `json:"gender"`
	ActivityLevel      string   `json:"activityLevel"`
	GoalType           string   `json:"goalType"`
	DietaryPreferences []string `json:"dietaryPreferences"`
	DislikedFoods      string   `json:"dislikedFoods"`
	Allergies          string   `json:"allergies"`
	DailyCalorieTarget *int     `json:"dailyCalorieTarget,omitempty"`
	TargetWaterMl      int      `json:"targetWaterMl"`
	TargetCoffeeCups   int      `json:"targetCoffeeCups"`
}

type UpdatePreferencesRequest struct {
	Preferences   []string       `json:"preferences"`
	Allergies     []AllergyInput `json:"allergies"`
	DislikedFoods []string       `json:"disliked_foods"`
}

type AllergyInput struct {
	Allergy  string  `json:"allergy"`
	Severity *string `json:"severity,omitempty"`
}

type PreferencesResponse struct {
	Preferences   []DietaryPreference `json:"preferences"`
	Allergies     []Allergy           `json:"allergies"`
	DislikedFoods []DislikedFood      `json:"disliked_foods"`
}

func (prefs *PreferencesResponse) ToResource() *response.PreferencesResponse {
	var dtoPrefs response.PreferencesResponse

	if prefs != nil {
		dtoPrefs.Preferences = make([]response.DietaryPreference, len(prefs.Preferences))
		for i, p := range prefs.Preferences {
			dtoPrefs.Preferences[i] = response.DietaryPreference{
				ID:         p.ID,
				Preference: p.Preference,
				CreatedAt:  p.CreatedAt,
			}
		}

		dtoPrefs.Allergies = make([]response.Allergy, len(prefs.Allergies))
		for i, a := range prefs.Allergies {
			dtoPrefs.Allergies[i] = response.Allergy{
				ID:        a.ID,
				Allergy:   a.Allergy,
				Severity:  a.Severity,
				CreatedAt: a.CreatedAt,
			}
		}

		dtoPrefs.DislikedFoods = make([]response.DislikedFood, len(prefs.DislikedFoods))
		for i, df := range prefs.DislikedFoods {
			dtoPrefs.DislikedFoods[i] = response.DislikedFood{
				ID:        df.ID,
				FoodName:  df.FoodName,
				CreatedAt: df.CreatedAt,
			}
		}
	}

	return &dtoPrefs
}

func (profile *UserProfile) ToResource() *response.UserProfile {
	return &response.UserProfile{
		ID:            profile.ID,
		UserID:        profile.UserID,
		DateOfBirth:   profile.DateOfBirth,
		Gender:        profile.Gender,
		HeightCm:      profile.HeightCm,
		WeightKg:      profile.WeightKg,
		ActivityLevel:      profile.ActivityLevel,
		Goal:               profile.Goal,
		Age:                profile.Age,
		DisplayName:        profile.DisplayName,
		TargetWeightKg:     profile.TargetWeightKg,
		DailyCalorieTarget: profile.DailyCalorieTarget,
		TargetWaterMl:      profile.TargetWaterMl,
		TargetCoffeeCups:   profile.TargetCoffeeCups,
		CreatedAt:          profile.CreatedAt,
		UpdatedAt:     profile.UpdatedAt,
	}
}
