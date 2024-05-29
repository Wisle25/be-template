//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
)

// External

func InitializedService() *users.UserRepository {
	// Repository
	wire.Build(
		repository.NewUserRepositoryPG,
		database.DB,
		generator.NewUUIDGenerator,
	)

	return nil
}
