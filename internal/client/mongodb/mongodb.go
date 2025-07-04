package mongodb

import (
	"context"
	"log"

	"github.com/nuea/backend-golang-test/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB interface {
	GetCollection(name string) *mongo.Collection
}

type mongoDB struct {
	client  *mongo.Client
	mongodb *mongo.Database
	cfg     *config.MongoDBConfig
}

func (m *mongoDB) GetCollection(name string) *mongo.Collection {
	if m.mongodb == nil {
		return nil
	}

	if m.mongodb.Collection(name) == nil {
		if err := m.mongodb.CreateCollection(context.Background(), name); err != nil {
			log.Println("Unable to create the collection: ", err)
			return nil
		}
	}
	return m.mongodb.Collection(name)
}

func ProvideMongoDBClient(cfg *config.AppConfig) (MongoDB, func(), error) {
	opt := options.Client().ApplyURI(cfg.MongoDB.Host).
		SetAuth(options.Credential{
			Username: cfg.MongoDB.User,
			Password: cfg.MongoDB.Password,
		}).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetMaxPoolSize(cfg.MongoDB.MaxPoolSize).
		SetMinPoolSize(cfg.MongoDB.MinPoolSize).
		SetHeartbeatInterval(cfg.MongoDB.HeartbeatInterval)
	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return nil, func() {}, err
	}

	var mongodb *mongo.Database
	if cfg.MongoDB.DatabaseName != "-" {
		mongodb = client.Database(cfg.MongoDB.DatabaseName)
	}

	log.Println("Start connecting to MongoDB:", cfg.MongoDB.DatabaseName)

	return &mongoDB{
			cfg:     &cfg.MongoDB,
			client:  client,
			mongodb: mongodb,
		}, func() {
			_ = client.Disconnect(context.Background())
		}, nil
}
