package repository

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/repository/user"
)

type Repository struct {
	user.UserRepository
}

var RepositorySet = wire.NewSet(
	user.ProvideUserRepository,

	wire.Struct(new(Repository), "*"),
)
