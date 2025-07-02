package user

import (
	"context"
	"time"

	"github.com/gotidy/ptr"
	"github.com/nuea/backend-golang-test/internal/repository"
	"github.com/nuea/backend-golang-test/internal/repository/user"
	"github.com/nuea/backend-golang-test/internal/types"
	"github.com/nuea/backend-golang-test/internal/util"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	userv1.UnimplementedUserServiceServer
	userrepo user.UserRepository
}

func ProvideUserGRPCService(repo *repository.Repository) (userv1.UserServiceServer, error) {
	return &grpcService{
		userrepo: repo.UserRepository,
	}, nil
}

func (g *grpcService) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	var email types.Email
	var err error
	if req.Email != "" {
		email, err = types.NewEmail(req.Email)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	newuser := user.NewUser()
	newuser.Name = req.Name
	newuser.Email = email

	newuser.Password, err = types.NewHashString(req.Password).Hash()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.CreatedBy != nil {
		newuser.CreatedBy = req.CreatedBy
	}
	newuser.CreatedBy = req.CreatedBy

	if err := g.userrepo.InsertOne(ctx, newuser); err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{}, nil
}

func (g *grpcService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {

	user, err := g.userrepo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	data, err := mapGRPCUser(user)
	if err != nil {
		return nil, err
	}
	return &userv1.GetUserResponse{User: data}, nil
}

func (g *grpcService) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	f := &user.UserFilter{}
	if req.Name != nil {
		f.Name = *req.Name
	}
	if req.Email != nil {
		email, err := types.NewEmail(*req.Email)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		f.Email = email
	}
	users, err := g.userrepo.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	datas, err := util.MapToSlice(mapGRPCUser, users)
	if err != nil {
		return nil, err
	}
	return &userv1.GetUsersResponse{
		Data: datas,
	}, nil
}

func (g *grpcService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	user, err := g.userrepo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		email, err := types.NewEmail(*req.Email)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		user.Email = email
	}
	user.UpdatedAt = time.Now().UTC()

	if err := g.userrepo.ReplaceOne(ctx, req.Id, user); err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{}, nil
}

func (g *grpcService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	user, err := g.userrepo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt = time.Now().UTC()
	user.DeletedAt = ptr.Time(time.Now().UTC())

	if err := g.userrepo.ReplaceOne(ctx, req.Id, user); err != nil {
		return nil, err
	}

	return &userv1.DeleteUserResponse{}, nil
}
