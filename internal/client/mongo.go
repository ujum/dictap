package client

import (
	"context"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type MongoClient struct {
	logger logger.Logger
	cfg    *config.MongoDatasourceConfig
	client *mongo.Client
}

func CreateMongoClient(cfg *config.MongoDatasourceConfig, log logger.Logger) (*MongoClient, error) {
	endpoint := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	opts := options.Client().ApplyURI("mongodb://" + endpoint)

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Debugf("connected to mongo %s", endpoint)

	return &MongoClient{
		logger: log,
		cfg:    cfg,
		client: client,
	}, nil
}

func (mongoClient *MongoClient) Query() {
	mongoClient.logger.Info("mongo client Query")
}

func (mongoClient *MongoClient) Disconnect() {
	if err := mongoClient.client.Disconnect(context.Background()); err != nil {
		mongoClient.logger.Error(err)
	}
	mongoClient.logger.Debug("mongo client disconnected")
}
