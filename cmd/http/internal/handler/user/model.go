package user

import "time"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type GetUsersRequest struct {
	Name  *string `form:"name,omitempty"`
	Email *string `form:"email,omitempty"`
}

type GetUsersResponse struct {
	Data []*User `json:"data"`
}

type GetUserResponse struct {
	User
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"emai,omitempty"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedBy *string    `json:"created_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
