//go:build wireinject
// +build wireinject

package progressservice

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"dietician.local/services/progress-service/internal"
	dailylogHandler "dietician.local/services/progress-service/internal/dailylog/handler"
	dailylogOrchestration "dietician.local/services/progress-service/internal/dailylog/orchestration"
	dailylogRepository "dietician.local/services/progress-service/internal/dailylog/repository"
	dailylogService "dietician.local/services/progress-service/internal/dailylog/service"
	"dietician.local/services/progress-service/internal/progress/handler"
	"dietician.local/services/progress-service/internal/progress/orchestration"
	"dietician.local/services/progress-service/internal/progress/repository"
	"dietician.local/services/progress-service/internal/progress/service"
	trackingHandler "dietician.local/services/progress-service/internal/tracking/handler"
	trackingOrchestration "dietician.local/services/progress-service/internal/tracking/orchestration"
	trackingRepository "dietician.local/services/progress-service/internal/tracking/repository"
	trackingService "dietician.local/services/progress-service/internal/tracking/service"
)

var progressSet = wire.NewSet(
	repository.NewProgressRepository,
	service.NewProgressService,
	orchestration.NewProgressOrchestrator,
	handler.NewProgressHandler,
)

var dailylogSet = wire.NewSet(
	dailylogRepository.NewDailyLogRepository,
	dailylogService.NewDailyLogService,
	dailylogOrchestration.NewDailyLogOrchestrator,
	dailylogHandler.NewDailyLogHandler,
)

var trackingSet = wire.NewSet(
	trackingRepository.NewTrackingRepository,
	trackingService.NewTrackingService,
	trackingOrchestration.NewTrackingOrchestrator,
	trackingHandler.NewTrackingHandler,
)

func InitRoute(db *sqlx.DB, logger *logrus.Logger) internal.IRoute {
	wire.Build(
		progressSet,
		dailylogSet,
		trackingSet,
		internal.NewRoute,
	)
	return nil
}
