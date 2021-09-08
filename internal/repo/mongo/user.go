package mongo

import (
	"context"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoMongo struct {
	collection *mongo.Collection
	log        logger.Logger
}

func NewUserRepoMongo(cfg *config.Config, log logger.Logger, collection *mongo.Collection) *UserRepoMongo {
	return &UserRepoMongo{
		collection: collection,
		log:        log,
	}
}

func (ur *UserRepoMongo) FindAll(ctx context.Context) ([]*domain.User, error) {
	cursor, err := ur.collection.Find(ctx, bson.D{})

	var users []*domain.User
	if err != nil {
		ur.log.Errorf("can't find all users, reason: %v", err)
		return users, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &users)
	return users, err
}

func (ur *UserRepoMongo) FindByUID(ctx context.Context, uid string) (*domain.User, error) {
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

func (ur *UserRepoMongo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	one := ur.collection.FindOne(ctx, bson.M{"email": email})
	if err := one.Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (ur *UserRepoMongo) Update(ctx context.Context, user *domain.User) error {
	user.ID = ""
	result, err := ur.collection.UpdateOne(ctx, bson.M{"uid": user.UID}, bson.M{"$set": user})
	if err != nil {
		ur.log.Errorf("can't update user, reason: %v", err)
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}
	return err
}

func (ur *UserRepoMongo) Create(ctx context.Context, user *domain.User) (string, error) {
	userID, err := create(ctx, ur.collection, user)
	if err == domain.ErrAlreadyExists {
		return userID, domain.ErrUserAlreadyExists
	}
	return userID, err
}

func (ur *UserRepoMongo) DeleteByUID(ctx context.Context, uid string) error {
	result, err := ur.collection.DeleteOne(ctx, bson.M{"uid": uid})
	if err != nil {
		ur.log.Errorf("can't delete user, reason: %v", err)
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
