package repo

import (
	"context"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	FindByUid(ctx context.Context, uid string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	FindAll(ctx context.Context) ([]*domain.User, error)
	DeleteByUid(ctx context.Context, uid string) error
	Update(ctx context.Context, user *domain.User) error
}

type UserRepoMongo struct {
	collection *mongo.Collection
}

func NewUserRepoMongo(cfg *config.Config, log logger.Logger, collection *mongo.Collection) *UserRepoMongo {
	return &UserRepoMongo{
		collection: collection,
	}
}

func (ur *UserRepoMongo) FindAll(ctx context.Context) ([]*domain.User, error) {
	cursor, err := ur.collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)

	var users []*domain.User
	if err != nil {
		return users, err
	}
	err = cursor.All(ctx, &users)
	return users, err
}

func (ur *UserRepoMongo) FindByUid(ctx context.Context, uid string) (*domain.User, error) {
	user := &domain.User{}
	one := ur.collection.FindOne(ctx, bson.M{"uid": uid})
	if err := one.Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (ur *UserRepoMongo) Update(ctx context.Context, user *domain.User) error {
	result, err := ur.collection.UpdateOne(ctx, bson.M{"uid": user.Uid}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}
	return err
}

func (ur *UserRepoMongo) Create(ctx context.Context, user *domain.User) error {
	_, err := ur.collection.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrUserAlreadyExists
	}
	return err
}

func (ur *UserRepoMongo) DeleteByUid(ctx context.Context, uid string) error {
	result, err := ur.collection.DeleteOne(ctx, bson.M{"uid": uid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
