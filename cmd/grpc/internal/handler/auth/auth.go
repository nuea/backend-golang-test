package auth

import (
	"context"

	"github.com/nuea/backend-golang-test/internal/repository"
	"github.com/nuea/backend-golang-test/internal/repository/user"
	"github.com/nuea/backend-golang-test/internal/types"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	userv1.UnimplementedAuthServiceServer
	userrepo user.UserRepository
}

func ProvideAuthGRPCService(repo *repository.Repository) (userv1.AuthServiceServer, error) {
	return &grpcService{
		userrepo: repo.UserRepository,
	}, nil
}

func (g *grpcService) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	email, err := types.NewEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := g.userrepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !types.NewHashString(user.Password).Equal(req.Password) {
		return nil, status.Error(codes.InvalidArgument, "Password is invalid.")
	}

	return &userv1.LoginResponse{
		UserId: user.ID.Hex(),
	}, nil
}
