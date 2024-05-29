//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
	"github.com/wisle25/be-template/infrastructures/security"
	"github.com/wisle25/be-template/infrastructures/validation"
)

func NewUserContainer() *users.UserRepository {
	// Repository
	wire.Build(
		repository.NewUserRepositoryPG,
		database.ProvideDB,
		generator.NewUUIDGenerator,
		use_case.NewAddUserUseCase,
		security.NewArgon2,
		validation.NewValidateUser,
		validation.NewValidator,
	)

	return nil
}
