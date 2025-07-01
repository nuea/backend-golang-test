//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	internalDI "github.com/nuea/backend-golang-test/internal/di"
)

var MainSet = wire.NewSet(
	internalDI.InternalSet,
	ProviderSet,

	wire.Struct(new(Container), "*"),
)

func InitContainer() (*Container, func(), error) {
	wire.Build(MainSet)

	return &Container{}, func() {}, nil
}
