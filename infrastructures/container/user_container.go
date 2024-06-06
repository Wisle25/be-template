//go:build wireinject
// +build wireinject

package container

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/cache"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
	"github.com/wisle25/be-template/infrastructures/security"
	"github.com/wisle25/be-template/infrastructures/validation"
)

func NewUserContainer(config *commons.Config, db *sql.DB, client *redis.Client) *use_case.UserUseCase {
	// Repository
	wire.Build(
		repository.NewUserRepositoryPG,
		cache.NewRedisCache,
		generator.NewUUIDGenerator,
		security.NewBcrypt,
		validation.NewValidateUser,
		validation.NewValidator,
		validation.NewValidatorTranslator,
		security.NewJwtToken,
		use_case.NewUserUseCase,
	)

	return nil
}
