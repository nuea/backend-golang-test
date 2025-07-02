package user

import (
	"context"

	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
)

type grpcService struct {
	userv1.UnimplementedUserServiceServer
}

func ProvideUserGRPCService() (userv1.UserServiceServer, error) {
	return &grpcService{}, nil
}

func (g *grpcService) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	return &userv1.CreateResponse{}, nil
}
