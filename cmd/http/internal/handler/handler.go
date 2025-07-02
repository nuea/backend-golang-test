package handler

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/cmd/http/internal/handler/user"
)

type Handlers struct {
	UserHandler *user.Handler
}

var HandlerSet = wire.NewSet(
	user.ProvideUserHandler,

	wire.Struct(new(Handlers), "*"),
)
