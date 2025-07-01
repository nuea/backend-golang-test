package di

import "github.com/nuea/backend-golang-test/cmd/grpc/internal/server"

type Container struct {
	server *server.GRPCServer
}

func (c *Container) Run() {
	c.server.Serve()
}
