package service

import (
	"context"
	"errors"
	"strings"

	"dietician.local/packages/constants"
	"dietician.local/packages/utils"
	"dietician.local/services/account-service/internal/profile/model"
	"dietician.local/services/account-service/internal/profile/repository"
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.ProfileNotFound))
	}
	return profile, nil
}

func (s *ProfileService) UpsertProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.UserProfile, error) {
	profile := &model.UserProfile{
		UserID:        userID,
		DateOfBirth:   req.DateOfBirth,
		Gender:        req.Gender,
		HeightCm:      req.HeightCm,
		ActivityLevel: req.ActivityLevel,
		Goal:          req.Goal,
	}
	if err := s.repo.Upsert(ctx, profile); err != nil {
		return nil, err
	}
	return s.repo.GetByUserID(ctx, userID)
}

func (s *ProfileService) Onboarding(ctx context.Context, userID string, req *model.OnboardingRequest) (*model.UserProfile, error) {
	var dislikedFoods []string
	if req.DislikedFoods != "" {
		parts := strings.Split(req.DislikedFoods, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				dislikedFoods = append(dislikedFoods, trimmed)
			}
		}
	}

	var allergies []model.AllergyInput
	if req.Allergies != "" {
		parts := strings.Split(req.Allergies, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				allergies = append(allergies, model.AllergyInput{
					Allergy: trimmed,
				})
			}
		}
	}

	profile := &model.UserProfile{
		UserID:             userID,
		Age:                &req.Age,
		DisplayName:        &req.Name,
		Gender:             &req.Gender,
		ActivityLevel:      &req.ActivityLevel,
		HeightCm:           &req.HeightCm,
		TargetWeightKg:     &req.TargetWeightKg,
		Goal:               &req.GoalType,
		DailyCalorieTarget: req.DailyCalorieTarget,
		TargetWaterMl:      &req.TargetWaterMl,
		TargetCoffeeCups:   &req.TargetCoffeeCups,
	}

	if err := s.repo.Upsert(ctx, profile); err != nil {
		return nil, err
	}

	if err := s.repo.ReplacePreferences(ctx, userID, req.DietaryPreferences); err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceDislikedFoods(ctx, userID, dislikedFoods); err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceAllergies(ctx, userID, allergies); err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(ctx, userID)
}

func (s *ProfileService) GetPreferences(ctx context.Context, userID string) (*model.PreferencesResponse, error) {
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	allergies, err := s.repo.GetAllergies(ctx, userID)
	if err != nil {
		return nil, err
	}
	foods, err := s.repo.GetDislikedFoods(ctx, userID)
	if err != nil {
		return nil, err
	}

	if prefs == nil {
		prefs = []model.DietaryPreference{}
	}
	if allergies == nil {
		allergies = []model.Allergy{}
	}
	if foods == nil {
		foods = []model.DislikedFood{}
	}

	return &model.PreferencesResponse{
		Preferences:   prefs,
		Allergies:     allergies,
		DislikedFoods: foods,
	}, nil
}

func (s *ProfileService) UpdatePreferences(ctx context.Context, userID string, req *model.UpdatePreferencesRequest) (*model.PreferencesResponse, error) {
	if err := s.repo.ReplacePreferences(ctx, userID, req.Preferences); err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceAllergies(ctx, userID, req.Allergies); err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceDislikedFoods(ctx, userID, req.DislikedFoods); err != nil {
		return nil, err
	}
	return s.GetPreferences(ctx, userID)
}
