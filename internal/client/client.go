package client

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/client/mongodb"
)

type Clients struct {
	MongoDB *mongodb.MongoDB
}

var ClientSet = wire.NewSet(
	mongodb.ProvideMongoDBClient,

	wire.Struct(new(Clients), "*"),
)
