package handler

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/cmd/http/internal/handler/auth"
	"github.com/nuea/backend-golang-test/cmd/http/internal/handler/user"
)

type Handlers struct {
	AuthHandler *auth.Handler
	UserHandler *user.Handler
}

var HandlerSet = wire.NewSet(
	auth.ProvideUserHandler,
	user.ProvideUserHandler,

	wire.Struct(new(Handlers), "*"),
)
