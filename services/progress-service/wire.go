//go:build wireinject
// +build wireinject

package progressservice

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"

	"dietician.local/services/progress-service/internal"
	"dietician.local/services/progress-service/internal/progress/handler"
	"dietician.local/services/progress-service/internal/progress/repository"
	"dietician.local/services/progress-service/internal/progress/service"
)

var progressSet = wire.NewSet(
	repository.NewProgressRepository,
	service.NewProgressService,
	handler.NewProgressHandler,
)

func InitRoute(db *sqlx.DB) internal.IRoute {
	wire.Build(
		progressSet,
		internal.NewRoute,
	)
	return nil
}
