package client

import (
	"context"
	"fmt"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type MongoClient struct {
	logger logger.Logger
	cfg    *config.MongoDatasourceConfig
	Client *mongo.Client
}

func CreateMongoClient(parentCtx context.Context, waitGroup *sync.WaitGroup, cfg *config.MongoDatasourceConfig, log logger.Logger) (*MongoClient, error) {
	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	opts := options.Client().ApplyURI("mongodb://" + endpoint)

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
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

	mc := &MongoClient{
		logger: log,
		cfg:    cfg,
		Client: client,
	}
	waitGroup.Add(1)
	go func() {
		<-parentCtx.Done()
		if discErr := mc.Disconnect(ctx); discErr != nil {
			log.Error(err)
		}
		defer waitGroup.Done()
	}()

	return mc, nil
}

func (mongoClient *MongoClient) Query() {
	mongoClient.logger.Info("mongo client Query")
}

func (mongoClient *MongoClient) Disconnect(ctx context.Context) error {
	if err := mongoClient.Client.Disconnect(ctx); err != nil {
		return err
	}
	mongoClient.logger.Debug("mongo client disconnected")
	return nil
}
