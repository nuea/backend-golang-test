package di

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/config"
)

var InternalSet = wire.NewSet(
	ConfigSet,
)

var ConfigSet = wire.NewSet(
	config.ProvideCofig,
)
