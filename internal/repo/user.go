package repo

import (
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
)

type UserRepo interface {
	GetByUid(uid string) string
}

type UserRepoMongo struct {
	client *client.MongoClient
}

func NewUserRepoMongo(cfg *config.Config, log logger.Logger, mongoClient *client.MongoClient) *UserRepoMongo {
	return &UserRepoMongo{
		client: mongoClient,
	}
}

func (ur *UserRepoMongo) GetByUid(uid string) string {
	ur.client.Query()
	return uid
}
