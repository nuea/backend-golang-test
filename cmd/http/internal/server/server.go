package server

import (
	"bytes"
	"context"
	"encoding/json"
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

type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
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

func WithRequestLoggerServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("HTTP request -", "method:", c.Request.Method, ", path:", c.Request.URL.Path)
		c.Next()
	}
}

func WithResponseLoggerServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapWriter := &ResponseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = wrapWriter
		c.Next()
		if errs := c.Errors.Last(); errs == nil {
			var body interface{}
			if err := json.Unmarshal(wrapWriter.body.Bytes(), &body); err != nil {
				log.Println("HTTP response -", "method:", c.Request.Method, ", path:", c.Request.URL.Path, ", http_status:", c.Writer.Status())
			} else {
				log.Println("HTTP response - ", "method", c.Request.Method, ", path", c.Request.URL.Path, ", http_status: ", c.Writer.Status(), ", response_body", body)
			}
		}
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
	sv.gin.Use(WithRequestLoggerServer())
	sv.gin.Use(WithResponseLoggerServer())
	sv.gin.Use(gin.Recovery())

	sv.load(h, m)

	return sv
}
