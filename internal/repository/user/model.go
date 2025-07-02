package user

import (
	"time"

	"github.com/nuea/backend-golang-test/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     types.Email        `bson:"email"`
	Password  string             `bson:"password"`
	CreatedBy *string            `bson:"created_by,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func NewUser() *User {
	return &User{
		CreatedBy: nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type UserFilter struct {
	User
}

func (f *UserFilter) Filter() bson.D {
	filter := bson.D{}
	if f.ID != primitive.NilObjectID {
		filter = append(filter, bson.E{Key: "_id", Value: f.ID})
	}
	if f.Name != "" {
		filter = append(filter, bson.E{Key: "name", Value: f.Name})
	}
	if f.Email != "" {
		filter = append(filter, bson.E{Key: "email", Value: f.Email})
	}
	return filter
}
