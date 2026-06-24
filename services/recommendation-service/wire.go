//go:build wireinject
// +build wireinject

package recoservice

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"

	"dietician.local/packages/openai"
	"dietician.local/services/recommendation-service/config"
	"dietician.local/services/recommendation-service/internal"
	
	dietplanHandler "dietician.local/services/recommendation-service/internal/dietplan/handler"
	dietplanOrchestrator "dietician.local/services/recommendation-service/internal/dietplan/orchestration"
	dietplanRepository "dietician.local/services/recommendation-service/internal/dietplan/repository"
	dietplanService "dietician.local/services/recommendation-service/internal/dietplan/service"

	"dietician.local/services/recommendation-service/internal/recommendation/handler"
	"dietician.local/services/recommendation-service/internal/recommendation/orchestration"
	"dietician.local/services/recommendation-service/internal/recommendation/service"
)

var recoSet = wire.NewSet(
	service.NewRecommendationService,
	orchestration.NewRecommendationOrchestrator,
	handler.NewRecommendationHandler,
)

var dietplanSet = wire.NewSet(
	dietplanRepository.NewDietPlanRepository,
	dietplanService.NewDietPlanService,
	dietplanOrchestrator.NewDietPlanOrchestrator,
	dietplanHandler.NewDietPlanHandler,
)

func InitRoute(cfg *config.RecommendationAppScheme, openaiService openai.Service, db *sqlx.DB) internal.IRoute {
	wire.Build(
		recoSet,
		dietplanSet,
		internal.NewRoute,
	)
	return nil
}
