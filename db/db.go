package db

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoOnce sync.Once
	ctx       context.Context
	clientErr error
	client    *mongo.Client
)

const (
	URI              = "mongodb://localhost:27017"
	DB               = "eltropy_db"
	ADMIN_COLLECTION = "admin"
)

func GetMongoClient() (*mongo.Client, error) {

	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(URI)
		client, clientErr = mongo.Connect(context.TODO(), clientOptions)
	})
	return client, clientErr
}
