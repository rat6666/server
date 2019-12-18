package client

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var optionsDB = options.Client().ApplyURI(getMongoURL())

//FactoryFunc ...
type FactoryFunc func() (Client, error)

//NewClient ...
func NewClient() (Client, error) {
	client, err := mongo.NewClient(optionsDB)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getMongoURL() string {
	if value, exists := os.LookupEnv("DATABASE_URL"); exists {
		return value
	}
	return "mongodb://localhost:27017"
}
