package di

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/config"
	"github.com/nuea/backend-golang-test/internal/repository"
)

var InternalSet = wire.NewSet(
	ConfigSet,
	client.ClientSet,
	repository.RepositorySet,
)

var ConfigSet = wire.NewSet(
	config.ProvideCofig,
)
