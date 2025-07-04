package user

import (
	"github.com/nuea/backend-golang-test/internal/repository/user"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapGRPCUser(user *user.User) (*userv1.User, error) {
	if user == nil {
		return nil, nil
	}

	response := &userv1.User{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     string(user.Email),
		CreatedBy: user.CreatedBy,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	if user.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*user.DeletedAt)
	}

	return response, nil
}
