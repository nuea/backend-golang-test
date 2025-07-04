package user

import (
	"context"
	"errors"
	"strings"

	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	InsertOne(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (user *User, err error)
	FindByEmail(ctx context.Context, email types.Email) (user *User, err error)
	Find(ctx context.Context, filter *UserFilter) (users []*User, err error)
	ReplaceOne(ctx context.Context, id string, user *User) error
	Count(ctx context.Context) (int64, error)
}

type repository struct {
	collection *mongo.Collection
}

func ProvideUserRepository(c *client.Clients) UserRepository {
	collection := c.MongoDB.GetCollection("user")
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	return &repository{
		collection: collection,
	}
}

func (r *repository) InsertOne(ctx context.Context, user *User) error {
	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) && strings.Contains(err.Error(), "email") {
			return errors.New("email already exists")
		}
		return err
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (user *User, err error) {
	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objid}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, err
}

func (r *repository) FindByEmail(ctx context.Context, email types.Email) (user *User, err error) {
	err = r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, err
}

func (r *repository) Find(ctx context.Context, filter *UserFilter) (users []*User, err error) {
	cur, err := r.collection.Find(ctx, filter.Filter())
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) ReplaceOne(ctx context.Context, id string, user *User) error {
	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if _, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objid}, user); err != nil {
		return err
	}
	return nil
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"deleted_at": nil})
}
