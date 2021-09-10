package client

import (
	"context"
	"fmt"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	logger logger.Logger
	cfg    *config.MongoDatasourceConfig
	Client *mongo.Client
}

func CreateMongoClient(parentCtx context.Context, cfg *config.MongoDatasourceConfig, log logger.Logger) (*MongoClient, error) {
	var endpoint string
	if cfg.Username != "" {
		endpoint = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	} else {
		endpoint = fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	}
	opts := options.Client().ApplyURI(endpoint)

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Debugf("connected to mongo [host: %s, port: %d]", cfg.Host, cfg.Port)

	mc := &MongoClient{
		logger: log,
		cfg:    cfg,
		Client: client,
	}
	return mc, nil
}

func (mongoClient *MongoClient) Disconnect(ctx context.Context) error {
	if err := mongoClient.Client.Disconnect(ctx); err != nil {
		return err
	}
	mongoClient.logger.Debug("mongo client disconnected")
	return nil
}
