//go:build wireinject
// +build wireinject

package accountservice

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"dietician.local/packages/tokenizer"
	"dietician.local/services/account-service/config"
	"dietician.local/services/account-service/internal"
	authrepository "dietician.local/services/account-service/internal/auth/repository"
	authservice "dietician.local/services/account-service/internal/auth/service"
	"dietician.local/services/account-service/internal/handler"
	profilerepository "dietician.local/services/account-service/internal/profile/repository"
	profileservice "dietician.local/services/account-service/internal/profile/service"
)

// authSet wires authentication dependencies
var authSet = wire.NewSet(
	authrepository.NewUserRepository,
	authrepository.NewRefreshTokenRepository,
	authrepository.NewOTPRepository,
	authservice.NewOTPService,
	authservice.NewRefreshTokenService,
	authservice.NewUserService,
	handler.NewAuthHandler,
)

// profileSet wires profile dependencies
var profileSet = wire.NewSet(
	profilerepository.NewProfileRepository,
	profileservice.NewProfileService,
	handler.NewProfileHandler,
)

// InitializeRoute creates the main Route instance with all dependencies wired.
func InitRoute(db *sqlx.DB, cfg *config.AccountAppScheme, sender authservice.EmailSender, rdb *redis.Client, tok tokenizer.ITokenizer) internal.IRoute {
	wire.Build(
		authSet,
		profileSet,
		internal.NewRoute,
	)
	return nil
}
