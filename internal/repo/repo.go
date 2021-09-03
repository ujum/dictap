package repo

import (
	"context"
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/repo/mongo"
	"github.com/ujum/dictap/pkg/logger"
)

type Repositories struct {
	UserRepo      UserRepo
	WordRepo      WordRepo
	WordGroupRepo WordGroupRepo
}

type UserRepo interface {
	FindByUid(ctx context.Context, uid string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (string, error)
	FindAll(ctx context.Context) ([]*domain.User, error)
	DeleteByUid(ctx context.Context, uid string) error
	Update(ctx context.Context, user *domain.User) error
}

type WordRepo interface {
	Create(ctx context.Context, word *domain.Word) (string, error)
	FindByGroup(ctx context.Context, groupID string) ([]*domain.Word, error)
	FindByName(ctx context.Context, name string) (*domain.Word, error)
	AddToGroup(ctx context.Context, wordID string, groupID string) error
	FindByNameAndGroup(ctx context.Context, wordName string, groupID string) (*domain.Word, error)
	RemoveFromGroup(ctx context.Context, name string, groupID string) error
}

type WordGroupRepo interface {
	Create(ctx context.Context, word *domain.WordGroup) (string, error)
	FindByIDAndUser(ctx context.Context, groupID string, userID string) (*domain.WordGroup, error)
	FindByLangAndUser(ctx context.Context, langID string, userID string, def bool) (*domain.WordGroup, error)
	FindAllByLangAndUser(ctx context.Context, langID string, userID string) ([]*domain.WordGroup, error)
}

func New(cfg *config.Config, log logger.Logger, clients *client.Clients) *Repositories {
	mongoDatabase := clients.Mongo.Client.Database(cfg.Datasource.Mongo.Schema)
	return &Repositories{
		UserRepo:      mongo.NewUserRepoMongo(cfg, log, mongoDatabase.Collection("users")),
		WordRepo:      mongo.NewWordRepoMongo(log, mongoDatabase.Collection("words")),
		WordGroupRepo: mongo.NewWordGroupRepoMongo(log, mongoDatabase.Collection("groups")),
	}
}
