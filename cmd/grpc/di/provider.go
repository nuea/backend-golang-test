package di

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/cmd/grpc/internal/server"
)

var ProviderSet = wire.NewSet(
	server.ProvideGRPCServer,
)
