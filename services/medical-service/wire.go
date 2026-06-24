//go:build wireinject
// +build wireinject

package medicalservice

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"dietician.local/services/medical-service/config"
	"dietician.local/services/medical-service/internal"
	"dietician.local/services/medical-service/internal/storage"
	"dietician.local/services/medical-service/internal/uploads/handler"
	"dietician.local/services/medical-service/internal/uploads/orchestration"
	"dietician.local/services/medical-service/internal/uploads/repository"
	"dietician.local/services/medical-service/internal/uploads/service"
)

var medicalSet = wire.NewSet(
	repository.NewMedicalRepository,
	service.NewMedicalService,
	orchestration.NewMedicalOrchestrator,
	handler.NewMedicalHandler,
)

func InitRoute(db *sqlx.DB, cfg *config.MedicalAppScheme, provider storage.Provider, logger *logrus.Logger) internal.IRoute {
	wire.Build(
		medicalSet,
		internal.NewRoute,
	)
	return nil
}
