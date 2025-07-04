package user

import (
	"context"
	"errors"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/nuea/backend-golang-test/internal/repository"
	"github.com/nuea/backend-golang-test/internal/repository/user"
	"github.com/nuea/backend-golang-test/internal/types"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUserRepository struct {
	mock.Mock
	user.UserRepository
}

func (m *mockUserRepository) InsertOne(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) Find(ctx context.Context, filter *user.UserFilter) ([]*user.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *mockUserRepository) ReplaceOne(ctx context.Context, id string, u *user.User) error {
	args := m.Called(ctx, id, u)
	return args.Error(0)
}

func TestProvideUserGRPCService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv, err := ProvideUserGRPCService(&repository.Repository{UserRepository: repo})

		assert.NotNil(t, sv)
		assert.NoError(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.CreateUserRequest{
			Name:      "test",
			Email:     "test@example.com",
			Password:  "password",
			CreatedBy: ptr.String("service"),
		}

		repo.On("InsertOne", ctx, mock.Anything).Return(nil).Once()

		res, err := sv.CreateUser(ctx, req)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.CreateUserRequest{
			Name:     "test",
			Email:    "invalid email",
			Password: "password",
		}
		msgerr := "mail: no angle-addr"
		res, err := sv.CreateUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), msgerr)
	})

	t.Run("invalid password", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.CreateUserRequest{
			Name:     "test",
			Email:    "test@example.com",
			Password: "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
		}

		res, err := sv.CreateUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Equal(t, "bcrypt: password length exceeds 72 bytes", st.Message())
	})

	t.Run("other error", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.CreateUserRequest{
			Name:     "test",
			Email:    "test@example.com",
			Password: "password",
		}

		msgerr := errors.New("internal server error")
		repo.On("InsertOne", ctx, mock.Anything).Return(msgerr).Once()

		res, err := sv.CreateUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr, err)
		repo.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	uid := primitive.NewObjectID()
	req := &userv1.GetUserRequest{
		Id: uid.Hex(),
	}
	muser := &user.User{
		ID:       uid,
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}

		repo.On("FindByID", ctx, req.Id).Return(muser, nil).Once()

		res, err := sv.GetUser(ctx, req)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.Equal(t, req.Id, res.User.Id)
		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		msgerr := errors.New("user not found")

		repo.On("FindByID", ctx, mock.Anything).Return(nil, msgerr).Once()

		res, err := sv.GetUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr, err)
		repo.AssertExpectations(t)
	})
}

func TestGetUsers(t *testing.T) {
	ctx := context.Background()
	req := &userv1.GetUsersRequest{
		Name:  ptr.String("test"),
		Email: ptr.String("test@example.com"),
	}
	expuser := []*user.User{
		{ID: primitive.NewObjectID(), Name: "test", Email: "test@example.com", Password: "password"},
		{ID: primitive.NewObjectID(), Name: "test", Email: "testtest@example.com", Password: "password"},
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		f := &user.UserFilter{}
		f.Name = *req.Name
		f.Email = types.Email(*req.Email)

		repo.On("Find", ctx, f).Return(expuser, nil).Once()

		res, err := sv.GetUsers(ctx, req)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.Equal(t, len(expuser), len(res.Data))
		assert.Equal(t, expuser[0].ID.Hex(), res.Data[0].Id)
		assert.Equal(t, expuser[1].ID.Hex(), res.Data[1].Id)
		repo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req.Email = ptr.String("invalid email")

		msgerr := "mail: no angle-addr"
		res, err := sv.GetUsers(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), msgerr)
	})

	t.Run("other error", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req.Email = ptr.String("testtest@example.com")

		msgerr := errors.New("internal server error")
		repo.On("Find", ctx, mock.Anything).Return(nil, msgerr).Once()

		res, err := sv.GetUsers(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr.Error(), err.Error())
		repo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	req := &userv1.UpdateUserRequest{
		Id:    primitive.NewObjectID().Hex(),
		Name:  ptr.String("test"),
		Email: ptr.String("test@example.com"),
	}
	muser := &user.User{
		ID:       primitive.NewObjectID(),
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}

		repo.On("FindByID", ctx, req.Id).Return(muser, nil).Once()
		repo.On("ReplaceOne", ctx, req.Id, muser).Return(nil).Once()

		res, err := sv.UpdateUser(ctx, req)
		assert.NotNil(t, res)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		msgerr := errors.New("user not found")

		repo.On("FindByID", ctx, mock.Anything).Return(nil, msgerr).Once()

		res, err := sv.UpdateUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr, err)
		repo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req.Email = ptr.String("invalid email")

		repo.On("FindByID", ctx, req.Id).Return(muser, nil).Once()

		msgerr := "mail: no angle-addr"
		res, err := sv.UpdateUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), msgerr)
		repo.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req.Email = ptr.String("testtest@example.com")
		msgerr := errors.New("internal server error")

		repo.On("FindByID", ctx, req.Id).Return(muser, nil).Once()
		repo.On("ReplaceOne", ctx, req.Id, muser).Return(msgerr).Once()
		res, err := sv.UpdateUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr.Error(), err.Error())
		repo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	uid := primitive.NewObjectID()
	req := &userv1.DeleteUserRequest{
		Id: uid.Hex(),
	}
	muser := &user.User{
		ID:       uid,
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}

		repo.On("FindByID", ctx, req.Id).Return(muser, nil).Once()
		repo.On("ReplaceOne", ctx, req.Id, muser).Return(nil).Once()

		res, err := sv.DeleteUser(ctx, req)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		msgerr := errors.New("user not found")

		repo.On("FindByID", ctx, mock.Anything).Return(nil, msgerr).Once()

		res, err := sv.DeleteUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr, err)
		repo.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		msgerr := errors.New("internal server error")

		repo.On("FindByID", ctx, req.Id).Return(muser, nil).Once()
		repo.On("ReplaceOne", ctx, req.Id, muser).Return(msgerr).Once()

		res, err := sv.DeleteUser(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Equal(t, msgerr, err)
		repo.AssertExpectations(t)
	})
}
