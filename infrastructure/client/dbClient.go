package client

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client ...
type Client interface {
	Disconnect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

type client struct{}

func (c *client) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	database := c.Database(name)
	return database
}

func (c *client) Disconnect(ctx context.Context) error {
	err := c.Disconnect(ctx)
	return err
}
