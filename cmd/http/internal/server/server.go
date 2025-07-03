package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/cmd/http/internal/handler"
	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/config"
	"github.com/nuea/backend-golang-test/internal/middleware"
	"github.com/oklog/run"
)

type HTTPServer struct {
	cfg    *config.AppConfig
	gin    *gin.Engine
	srv    *http.Server
	client *client.Clients
}

func (s *HTTPServer) Serve() {
	g := &run.Group{}
	g.Add(func() error {
		s.srv = &http.Server{
			Addr:    fmt.Sprintf(":%s", s.cfg.System.HTTPPort),
			Handler: s.gin.Handler(),
		}
		log.Println("HTTP Server - started at ip address", s.srv.Addr)
		return s.srv.ListenAndServe()
	}, func(error) {
		if err := s.srv.Shutdown(context.Background()); err != nil {
			log.Println("Failed to close HTTP server")
		}
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		log.Println("HTTP Server - failed")
		os.Exit(1)
	}
}

func (s *HTTPServer) load(h *handler.Handlers, m *middleware.Middleware) {
	registerRouter(s.gin, h, m)
}

func ProvideHTTPServer(cfg *config.AppConfig, h *handler.Handlers, c *client.Clients, m *middleware.Middleware) *HTTPServer {
	sv := &HTTPServer{
		cfg:    cfg,
		gin:    gin.New(),
		srv:    &http.Server{},
		client: c,
	}

	sv.load(h, m)

	sv.gin.Use(gin.Logger())
	sv.gin.Use(gin.Recovery())

	return sv
}
