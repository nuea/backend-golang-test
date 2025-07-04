package user

import (
	"context"
	"strings"
	"testing"

	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestInsertOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	user := &User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.InsertOne(context.Background(), user)

		assert.Nil(t, err)
	})

	mt.Run("duplicate key error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "E11000 duplicate key error collection: test.user index: email_1 dup key: { email: \"test@example.com\" }",
		}))

		err := repo.InsertOne(context.Background(), user)

		assert.EqualError(t, err, "email already exists")
	})

	mt.Run("other error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		msg := "some other write error"

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    1,
			Message: msg,
		}))

		err := repo.InsertOne(context.Background(), user)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Error(t, err, msg)
	})
}

type mockMongoDB struct {
	mt *mtest.T
}

func (m *mockMongoDB) GetCollection(name string) *mongo.Collection {
	return m.mt.Coll
}

func TestProvideUserRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mm := new(mockMongoDB)
		mm.mt = mt

		repo := ProvideUserRepository(&client.Clients{
			MongoDB: mm,
		})
		assert.NotNil(t, repo)
	})

}

func TestFindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	uid := primitive.NewObjectID()
	user := &User{
		ID:       uid,
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: user.ID},
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
				{Key: "password", Value: user.Password},
			}))
		res, err := repo.FindByID(context.Background(), user.ID.Hex())

		assert.Nil(t, err)
		assert.Equal(t, user, res)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.user", mtest.FirstBatch))

		user, err := repo.FindByID(context.Background(), uid.Hex())

		assert.Nil(t, user)
		assert.EqualError(t, err, "user not found")
	})

	mt.Run("invalid id", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		errmsg := "the provided hex string is not a valid ObjectID"

		user, err := repo.FindByID(context.Background(), "invalid id")

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Error(t, err, errmsg)

	})

	mt.Run("other error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		uid := primitive.NewObjectID()
		msg := "internal server error"

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 1, Message: msg}))

		user, err := repo.FindByID(context.Background(), uid.Hex())

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Error(t, err, msg)
	})
}

func TestFindByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	user := &User{
		ID:       primitive.NewObjectID(),
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: user.ID},
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
				{Key: "password", Value: user.Password},
			}))
		res, err := repo.FindByEmail(context.Background(), user.Email)
		assert.Nil(t, err)
		assert.Equal(t, user, res)
	})

	mt.Run("user not found", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.user", mtest.FirstBatch))

		user, err := repo.FindByEmail(context.Background(), user.Email)

		assert.Nil(t, user)
		assert.EqualError(t, err, "user not found")
	})

	mt.Run("other error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		msg := "internal server error"

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 1, Message: msg}))

		user, err := repo.FindByEmail(context.Background(), user.Email)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Error(t, err, msg)
	})
}

func TestFind(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	expuser := []*User{
		{ID: primitive.NewObjectID(), Name: "test", Email: "test@example.com", Password: "password"},
		{ID: primitive.NewObjectID(), Name: "test", Email: "testtest@example.com", Password: "password"},
	}

	curone := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch,
		bson.D{
			{Key: "_id", Value: expuser[0].ID},
			{Key: "name", Value: expuser[0].Name},
			{Key: "email", Value: expuser[0].Email},
			{Key: "password", Value: expuser[0].Password},
		})

	curtwo := mtest.CreateCursorResponse(2, "test.user", mtest.NextBatch,
		bson.D{
			{Key: "_id", Value: expuser[1].ID},
			{Key: "name", Value: expuser[1].Name},
			{Key: "email", Value: expuser[1].Email},
			{Key: "password", Value: expuser[1].Password},
		})

	endcur := mtest.CreateCursorResponse(0, "test.user", mtest.NextBatch)

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		f := &UserFilter{}

		mt.AddMockResponses(curone, curtwo, endcur)

		users, err := repo.Find(context.Background(), f)

		assert.Nil(t, err)
		assert.Equal(t, len(expuser), len(users))
		assert.Equal(t, expuser[0], users[0])
	})

	mt.Run("success with filter", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		f := &UserFilter{}
		f.Email = expuser[0].Email

		mt.AddMockResponses(curone, curtwo, endcur)

		users, err := repo.Find(context.Background(), f)

		assert.Nil(t, err)
		assert.Equal(t, len(expuser), len(users))
		assert.Equal(t, expuser[0], users[0])
	})

	mt.Run("other error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		msg := "internal server error"

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 1, Message: msg}))

		user, err := repo.Find(context.Background(), &UserFilter{})

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Error(t, err, msg)
	})

	mt.Run("cursor all error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		filter := &UserFilter{}
		msg := "cursor next error"

		curerr := mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 1, Message: msg})
		mt.AddMockResponses(curone, curerr)

		users, err := repo.Find(context.Background(), filter)

		assert.Nil(t, users)
		assert.NotNil(t, err)
		assert.Error(t, err, msg)
	})

}

func TestReplaceOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	user := &User{
		ID:       primitive.NewObjectID(),
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: user.ID},
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
				{Key: "password", Value: user.Password},
			}))
		err := repo.ReplaceOne(context.Background(), user.ID.Hex(), user)

		assert.Nil(t, err)
	})

	mt.Run("invalid id", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		errmsg := "the provided hex string is not a valid ObjectID"

		err := repo.ReplaceOne(context.Background(), "invalid id", user)

		assert.NotNil(t, err)
		assert.Error(t, err, errmsg)

	})

	mt.Run("other error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		msg := "internal server error"

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 1, Message: msg}))

		err := repo.ReplaceOne(context.Background(), user.ID.Hex(), user)
		if err != nil && !strings.Contains(err.Error(), msg) {
			mt.Fatalf("expected 'internal server error', got '%v'", err)
		}

		assert.NotNil(t, err)
		assert.Error(t, err, msg)
	})
}

func TestCount(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{{Key: "n", Value: 2}}))

		count, err := repo.Count(context.Background())

		assert.Nil(t, err)
		assert.Equal(mt, int64(2), count, "expected count 2, got %v", count)
	})

	mt.Run("other error", func(mt *mtest.T) {
		repo := &repository{collection: mt.Coll}
		msg := "cursor next error"

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 1, Message: msg}))

		_, err := repo.Count(context.Background())

		assert.NotNil(t, err)
		assert.Error(t, err, msg)
	})

}
