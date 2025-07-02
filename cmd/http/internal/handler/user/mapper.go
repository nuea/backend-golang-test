package user

import (
	"github.com/gotidy/ptr"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
)

func mapToUser(user *userv1.User) (*User, error) {
	response := &User{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}

	if user.DeletedAt != nil {
		response.DeletedAt = ptr.Of(user.DeletedAt.AsTime())
	}

	return response, nil
}
