//go:build wireinject
// +build wireinject

package recoservice

import (
	"github.com/google/wire"

	"dietician.local/packages/openai"
	"dietician.local/services/recommendation-service/config"
	"dietician.local/services/recommendation-service/internal"
	"dietician.local/services/recommendation-service/internal/recommendation/handler"
	"dietician.local/services/recommendation-service/internal/recommendation/service"
)

var recoSet = wire.NewSet(
	service.NewRecommendationService,
	handler.NewRecommendationHandler,
)

func InitRoute(cfg *config.RecommendationAppScheme, openaiService openai.Service) internal.IRoute {
	wire.Build(
		recoSet,
		internal.NewRoute,
	)
	return nil
}
