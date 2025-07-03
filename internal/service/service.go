package service

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/service/auth"
)

type Service struct {
	auth.AuthService
}

var ServiceSet = wire.NewSet(
	auth.ProvideAuthenticationService,

	wire.Struct(new(Service), "*"),
)
