//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
	"github.com/wisle25/be-template/infrastructures/security"
	"github.com/wisle25/be-template/infrastructures/validation"
)

func NewUserContainer() *use_case.UserUseCase {
	// Repository
	wire.Build(
		repository.NewUserRepositoryPG,
		database.ProvideDB,
		generator.NewUUIDGenerator,
		security.NewArgon2,
		validation.NewValidateUser,
		validation.NewValidator,
		validation.NewValidatorTranslator,
		use_case.NewAddUserUseCase,
	)

	return nil
}
