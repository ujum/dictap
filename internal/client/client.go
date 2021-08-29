package client

import (
	"context"
	"github.com/pkg/errors"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
	"sync"
)

type Clients struct {
	Mongo     *MongoClient
	waitGroup *sync.WaitGroup
}

func New(ctx context.Context, cfg *config.DatasourceConfig, appLogger logger.Logger) (*Clients, error) {
	waitGroup := &sync.WaitGroup{}
	clients := &Clients{
		waitGroup: waitGroup,
	}
	mongo, err := CreateMongoClient(ctx, waitGroup, cfg.Mongo, appLogger)
	if err != nil {
		return clients, errors.Wrap(err, "can't create mongo client")
	}
	clients.Mongo = mongo

	return clients, nil
}

func (clients *Clients) WaitDisconnect() {
	if clients.waitGroup != nil {
		clients.waitGroup.Wait()
	}
}
