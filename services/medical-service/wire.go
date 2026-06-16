//go:build wireinject
// +build wireinject

package medicalservice

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"

	"dietician.local/services/medical-service/config"
	"dietician.local/services/medical-service/internal"
	"dietician.local/services/medical-service/internal/medical/handler"
	"dietician.local/services/medical-service/internal/medical/repository"
	"dietician.local/services/medical-service/internal/medical/service"
	"dietician.local/services/medical-service/internal/storage"
)

var medicalSet = wire.NewSet(
	repository.NewMedicalRepository,
	service.NewMedicalService,
	handler.NewMedicalHandler,
)

func InitRoute(db *sqlx.DB, cfg *config.MedicalAppScheme, provider storage.Provider) internal.IRoute {
	wire.Build(
		medicalSet,
		internal.NewRoute,
	)
	return nil
}
