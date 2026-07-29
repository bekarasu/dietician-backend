package orchestration

import (
	"context"

	"dietician.local/services/account-service/internal/profile/model"
	"dietician.local/services/account-service/internal/profile/service"
)

type IProfileOrchestrator interface {
	GetProfile(ctx context.Context, userID string) (*model.UserProfile, error)
	UpsertProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.UserProfile, error)
	Onboarding(ctx context.Context, userID string, req *model.OnboardingRequest) (*model.UserProfile, error)
	GetPreferences(ctx context.Context, userID string) (*model.PreferencesResponse, error)
	UpdatePreferences(ctx context.Context, userID string, req *model.UpdatePreferencesRequest) (*model.PreferencesResponse, error)
}

type profileOrchestrator struct {
	profileService *service.ProfileService
}

func NewProfileOrchestrator(profileService *service.ProfileService) IProfileOrchestrator {
	return &profileOrchestrator{profileService: profileService}
}

func (o *profileOrchestrator) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	return o.profileService.GetProfile(ctx, userID)
}

func (o *profileOrchestrator) UpsertProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.UserProfile, error) {
	return o.profileService.UpsertProfile(ctx, userID, req)
}

func (o *profileOrchestrator) Onboarding(ctx context.Context, userID string, req *model.OnboardingRequest) (*model.UserProfile, error) {
	return o.profileService.Onboarding(ctx, userID, req)
}

func (o *profileOrchestrator) GetPreferences(ctx context.Context, userID string) (*model.PreferencesResponse, error) {
	return o.profileService.GetPreferences(ctx, userID)
}

func (o *profileOrchestrator) UpdatePreferences(ctx context.Context, userID string, req *model.UpdatePreferencesRequest) (*model.PreferencesResponse, error) {
	return o.profileService.UpdatePreferences(ctx, userID, req)
}
