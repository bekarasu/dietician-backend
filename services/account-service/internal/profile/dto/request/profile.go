package request

import "dietician.local/services/account-service/internal/profile/model"

type UpdateProfileRequest struct {
	DateOfBirth   *string  `json:"date_of_birth,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender        *string  `json:"gender,omitempty"`
	HeightCm      *float64 `json:"height_cm,omitempty" validate:"omitempty,gt=0"`
	WeightKg      *float64 `json:"weight_kg,omitempty" validate:"omitempty,gt=0"`
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
	Name               string   `json:"name" validate:"required"`
	Age                int      `json:"age" validate:"required,gt=0"`
	HeightCm           float64  `json:"heightCm" validate:"required,gt=0"`
	WeightKg           float64  `json:"weightKg" validate:"required,gt=0"`
	TargetWeightKg     float64  `json:"targetWeightKg" validate:"required,gt=0"`
	Gender             string   `json:"gender" validate:"required"`
	ActivityLevel      string   `json:"activityLevel" validate:"required"`
	GoalType           string   `json:"goalType" validate:"required"`
	DietaryPreferences []string `json:"dietaryPreferences" validate:"dive"`
	DislikedFoods      string   `json:"dislikedFoods"`
	Allergies          string   `json:"allergies"`
	DailyCalorieTarget *int     `json:"dailyCalorieTarget,omitempty"`
	TargetWaterMl      int      `json:"targetWaterMl" validate:"required,gt=0"`
	TargetCoffeeCups   int      `json:"targetCoffeeCups" validate:"required,gte=0"`
}
type UpdatePreferencesRequest struct {
	Preferences   []string       `json:"preferences" validate:"dive,required"`
	Allergies     []AllergyInput `json:"allergies" validate:"dive"`
	DislikedFoods []string       `json:"disliked_foods" validate:"dive,required"`
}

type AllergyInput struct {
	Allergy  string  `json:"allergy" validate:"required"`
	Severity *string `json:"severity,omitempty"`
}

func (req *UpdatePreferencesRequest) ToDomain() *model.UpdatePreferencesRequest {
	var allergies []model.AllergyInput
	for _, a := range req.Allergies {
		allergies = append(allergies, model.AllergyInput{
			Allergy:  a.Allergy,
			Severity: a.Severity,
		})
	}
	return &model.UpdatePreferencesRequest{
		Preferences:   req.Preferences,
		Allergies:     allergies,
		DislikedFoods: req.DislikedFoods,
	}
}

func (req *UpdateProfileRequest) ToDomain() *model.UpdateProfileRequest {
	return &model.UpdateProfileRequest{
		DateOfBirth:   req.DateOfBirth,
		Gender:        req.Gender,
		HeightCm:      req.HeightCm,
		WeightKg:      req.WeightKg,
		ActivityLevel:      req.ActivityLevel,
		Goal:               req.Goal,
		Age:                req.Age,
		DisplayName:        req.DisplayName,
		TargetWeightKg:     req.TargetWeightKg,
		DailyCalorieTarget: req.DailyCalorieTarget,
		TargetWaterMl:      req.TargetWaterMl,
		TargetCoffeeCups:   req.TargetCoffeeCups,
	}
}

func (req *OnboardingRequest) ToDomain() *model.OnboardingRequest {
	return &model.OnboardingRequest{
		Name:               req.Name,
		Age:                req.Age,
		HeightCm:           req.HeightCm,
		WeightKg:           req.WeightKg,
		TargetWeightKg:     req.TargetWeightKg,
		Gender:             req.Gender,
		ActivityLevel:      req.ActivityLevel,
		GoalType:           req.GoalType,
		DietaryPreferences: req.DietaryPreferences,
		DislikedFoods:      req.DislikedFoods,
		Allergies:          req.Allergies,
		DailyCalorieTarget: req.DailyCalorieTarget,
		TargetWaterMl:      req.TargetWaterMl,
		TargetCoffeeCups:   req.TargetCoffeeCups,
	}
}
