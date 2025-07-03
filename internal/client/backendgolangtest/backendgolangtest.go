package backendgolangtest

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/nuea/backend-golang-test/internal/config"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BackendGolangTestGRPCService struct {
	userv1.UserServiceClient
	userv1.AuthServiceClient
}

type APIClient struct {
	conn *grpc.ClientConn
}

func NewDefaultGRPCClient(target string, du time.Duration, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if opts == nil {
		opts = make([]grpc.DialOption, 0)
	}

	baseOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainStreamInterceptor(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithIdleTimeout(du),
		grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, du)
			defer cancelFunc()
			return invoker(ctxWithTimeout, method, req, reply, cc, opts...)
		}),
	}

	return grpc.NewClient(target, append(baseOpts, opts...)...)
}

func WithRequestLoggerUnaryClient() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("GRPC request client - ", "request:", req, "method:", method)
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}

func ProvideBackendGolangTestServiceGRPC(cfg *config.AppConfig) *APIClient {
	conn, err := NewDefaultGRPCClient(cfg.BackendGoTest.GRPCTarget, cfg.BackendGoTest.RequestTimeout, WithRequestLoggerUnaryClient())
	if err != nil {
		panic(err)
	}

	return &APIClient{
		conn: conn,
	}
}

func ProvideUserServiceClient(client *APIClient) userv1.UserServiceClient {
	return userv1.NewUserServiceClient(client.conn)
}

func ProvideAuthServiceClient(client *APIClient) userv1.AuthServiceClient {
	return userv1.NewAuthServiceClient(client.conn)
}
