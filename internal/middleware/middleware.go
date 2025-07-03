package middleware

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/internal/middleware/auth"
)

type Middleware struct {
	Auth auth.AuthMiddleware
}

var MiddlewareSet = wire.NewSet(
	auth.ProvideAuthMiddleware,

	wire.Struct(new(Middleware), "*"),
)
