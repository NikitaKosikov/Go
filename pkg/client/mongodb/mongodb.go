package mongodb

import (
	"context"
	"fmt"
	"test/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	timeout = 10 * time.Second
)

func NewClient(mc config.MongodbConfig) (db *mongo.Client, err error) {
	clientOptions := options.Client().ApplyURI(mc.URI)
	if mc.Username != "" && mc.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: mc.Username, Password: mc.Password,
		})
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mongo client, err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start background to mongo client, err: %v", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil,  fmt.Errorf("failed to connect mongo, err: %v", err)
	}

	return client, nil
}
