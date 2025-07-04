package client

import (
	"github.com/google/wire"
	begot "github.com/nuea/backend-golang-test/internal/client/backendgolangtest"
	"github.com/nuea/backend-golang-test/internal/client/mongodb"
)

type GRPCClients struct {
	*begot.BackendGolangTestGRPCService
}

type Clients struct {
	MongoDB mongodb.MongoDB
}

var ClientSet = wire.NewSet(
	mongodb.ProvideMongoDBClient,
	begot.ProvideBackendGolangTestServiceGRPC,
	begot.ProvideUserServiceClient,
	begot.ProvideAuthServiceClient,

	wire.Struct(new(begot.BackendGolangTestGRPCService), "*"),
	wire.Struct(new(GRPCClients), "*"),
	wire.Struct(new(Clients), "*"),
)
