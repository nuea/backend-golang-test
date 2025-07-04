package auth

import (
	"context"
	"errors"
	"testing"

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

func (m *mockUserRepository) FindByEmail(ctx context.Context, email types.Email) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func TestProvideAuthGRPCService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv, err := ProvideAuthGRPCService(&repository.Repository{UserRepository: repo})

		assert.NotNil(t, sv)
		assert.NoError(t, err)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	pwd := "password"
	hash, _ := types.NewHashString(pwd).Hash()
	uid := primitive.NewObjectID()

	muser := &user.User{
		ID:       uid,
		Email:    types.Email(pwd),
		Password: hash,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.LoginRequest{Email: email, Password: pwd}

		repo.On("FindByEmail", ctx, types.Email(email)).Return(muser, nil).Once()

		res, err := sv.Login(ctx, req)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.Equal(t, uid.Hex(), res.UserId)
		repo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.LoginRequest{Email: "invalid email", Password: pwd}
		msgerr := "mail: no angle-addr"
		res, err := sv.Login(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), msgerr)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.LoginRequest{Email: email, Password: pwd}
		msgerr := errors.New("user not found")

		repo.On("FindByEmail", ctx, types.Email(email)).Return(nil, msgerr).Once()

		res, err := sv.Login(ctx, req)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, msgerr)
		repo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		repo := new(mockUserRepository)
		sv := &grpcService{userrepo: repo}
		req := &userv1.LoginRequest{Email: email, Password: "wrong-password"}

		repo.On("FindByEmail", ctx, types.Email(email)).Return(muser, nil).Once()

		res, err := sv.Login(ctx, req)

		assert.Nil(t, res)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Equal(t, "Password is invalid.", st.Message())
		repo.AssertExpectations(t)
	})
}
