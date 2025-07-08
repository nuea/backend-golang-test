package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/nuea/backend-golang-test/cmd/grpc/internal/handler"
	"github.com/nuea/backend-golang-test/internal/config"
	"github.com/nuea/backend-golang-test/internal/repository"
	"github.com/nuea/backend-golang-test/internal/repository/user"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	cfg      *config.AppConfig
	srv      *grpc.Server
	userrepo user.UserRepository
}

func (s *GRPCServer) Serve() {
	g := &run.Group{}
	g.Add(func() error {
		ipaddr := fmt.Sprintf(":%s", s.cfg.GRPCConfig.GRPCPort)
		lis, err := net.Listen("tcp", ipaddr)
		if err != nil {
			panic(err)
		}

		log.Println("GRPC Server - started at ip address", ipaddr)
		s.srv.Serve(lis)
		return nil
	}, func(err error) {
		s.srv.GracefulStop()
		s.srv.Stop()
	})

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(func() error {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				{
					count, err := s.userrepo.Count(context.Background())
					if err != nil {
						log.Printf("Error counting users: %v", err)
					}
					log.Printf("Number of users in database: %d user", count)
				}
			case <-ctx.Done():
				return nil
			}
		}
	}, func(err error) {
		cancel()
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := g.Run(); err != nil {
		log.Println("GRPC Server - failed")
		os.Exit(1)
	}
}

func ProvideGRPCServer(cfg *config.AppConfig, h *handler.GrpcServices, r *repository.Repository) *GRPCServer {
	opt := make([]grpc.ServerOption, 0)
	opt = append(opt, grpc.Creds(insecure.NewCredentials()))
	opt = append(opt, grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    2 * time.Hour,
		Timeout: 20 * time.Second,
	}))

	s := &GRPCServer{
		cfg:      cfg,
		srv:      grpc.NewServer(opt...),
		userrepo: r.UserRepository,
	}

	handler.RegisterGrpcServices(s.srv, h)
	if !cfg.GRPCConfig.GRPCHealthcheckDisabled {
		healthgrpc.RegisterHealthServer(s.srv, health.NewServer())
	}

	if cfg.GRPCConfig.GRPCReflectionEnabled {
		reflection.Register(s.srv)
	}

	return s
}
