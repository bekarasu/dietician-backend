package request

import "dietician.local/services/account-service/internal/profile/model"

type UpdateProfileRequest struct {
	DateOfBirth   *string  `json:"date_of_birth,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender        *string  `json:"gender,omitempty"`
	HeightCm      *float64 `json:"height_cm,omitempty" validate:"omitempty,gt=0"`
	WeightKg      *float64 `json:"weight_kg,omitempty" validate:"omitempty,gt=0"`
	ActivityLevel *string  `json:"activity_level,omitempty"`
	Goal          *string  `json:"goal,omitempty"`
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
		ActivityLevel: req.ActivityLevel,
		Goal:          req.Goal,
	}
}
