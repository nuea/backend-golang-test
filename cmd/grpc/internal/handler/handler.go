package handler

import (
	"github.com/google/wire"
	"github.com/nuea/backend-golang-test/cmd/grpc/internal/handler/auth"
	"github.com/nuea/backend-golang-test/cmd/grpc/internal/handler/user"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"google.golang.org/grpc"
)

type GrpcServices struct {
	userv1.UserServiceServer
	userv1.AuthServiceServer
}

func RegisterGrpcServices(sv *grpc.Server, h *GrpcServices) {
	userv1.RegisterUserServiceServer(sv, h)
	userv1.RegisterAuthServiceServer(sv, h)
}

var HandlerSet = wire.NewSet(
	user.ProvideUserGRPCService,
	auth.ProvideAuthGRPCService,

	wire.Struct(new(GrpcServices), "*"),
)
