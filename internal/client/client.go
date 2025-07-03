package client

import (
	"github.com/google/wire"
	begot "github.com/nuea/backend-golang-test/internal/client/backendgolangtest"
	"github.com/nuea/backend-golang-test/internal/client/mongodb"
)

type Clients struct {
	MongoDB *mongodb.MongoDB
	*begot.BackendGolangTestGRPCService
}

var ClientSet = wire.NewSet(
	mongodb.ProvideMongoDBClient,
	begot.ProvideBackendGolangTestServiceGRPC,
	begot.ProvideUserServiceClient,
	begot.ProvideAuthServiceClient,

	wire.Struct(new(begot.BackendGolangTestGRPCService), "*"),
	wire.Struct(new(Clients), "*"),
)
