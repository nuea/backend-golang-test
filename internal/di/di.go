package di

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/config"
	"github.com/nuea/backend-golang-test/internal/middleware"
	"github.com/nuea/backend-golang-test/internal/repository"
	"github.com/nuea/backend-golang-test/internal/service"
)

var InternalSet = wire.NewSet(
	ConfigSet,
	client.ClientSet,
	repository.RepositorySet,
	service.ServiceSet,
	middleware.MiddlewareSet,
)

var ConfigSet = wire.NewSet(
	config.ProvideCofig,
)
