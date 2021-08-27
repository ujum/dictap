package client

import (
	"github.com/pkg/errors"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
)

type Clients struct {
	Mongo *MongoClient
}

func New(cfg *config.DatasourceConfig, appLogger logger.Logger) (*Clients, error) {
	client, err := CreateMongoClient(cfg.Mongo, appLogger)
	if err != nil {
		return nil, errors.Wrap(err, "can't create mongo client")
	}
	return &Clients{
		Mongo: client,
	}, nil
}
