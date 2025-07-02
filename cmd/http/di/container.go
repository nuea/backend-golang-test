package di

import (
	"github.com/nuea/backend-golang-test/cmd/http/internal/server"
)

type Container struct {
	server *server.HTTPServer
}

func (c *Container) Run() {
	c.server.Serve()
}
