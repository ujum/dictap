package client

import (
	"context"
	"github.com/pkg/errors"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
)

type Clients struct {
	logger logger.Logger
	Mongo  *MongoClient
}

func New(ctx context.Context, cfg *config.DatasourceConfig, appLogger logger.Logger) (*Clients, error) {
	clients := &Clients{
		logger: appLogger,
	}
	mongo, err := CreateMongoClient(ctx, cfg.Mongo, appLogger)
	if err != nil {
		return clients, errors.Wrap(err, "can't create mongo client")
	}
	clients.Mongo = mongo
	return clients, nil
}

func (clients *Clients) Disconnect() {
	if clients.Mongo != nil {
		if err := clients.Mongo.Disconnect(context.Background()); err != nil {
			clients.logger.Errorf("error during disconnecting mongo client: %v", err)
		}
	}
}
