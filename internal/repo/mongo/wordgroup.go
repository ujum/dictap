package mongo

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WordGroupRepoMongo struct {
	log        logger.Logger
	collection *mongo.Collection
}

func (wgr *WordGroupRepoMongo) FindAllByLangAndUser(ctx context.Context, langID string, userID string) ([]*domain.WordGroup, error) {
	userIDHex, _ := primitive.ObjectIDFromHex(userID)
	langIDHex, _ := primitive.ObjectIDFromHex(langID)
	var wgs []*domain.WordGroup

	cursor, err := wgr.collection.Find(ctx, bson.M{"user_id": userIDHex, "lang_id": langIDHex})
	if err != nil {
		wgr.log.Errorf("can't find word groups by user, reason: %v", err)
		return wgs, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &wgs); err != nil {
		return nil, err
	}
	return wgs, nil
}

func (wgr *WordGroupRepoMongo) FindByIDAndUser(ctx context.Context, groupID string, userID string) (*domain.WordGroup, error) {
	groupIDHEx, _ := primitive.ObjectIDFromHex(groupID)
	userIDHEx, _ := primitive.ObjectIDFromHex(userID)
	wg := &domain.WordGroup{}

	result := wgr.collection.FindOne(ctx, bson.M{"_id": groupIDHEx, "user_id": userIDHEx})
	if err := result.Decode(wg); err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrNotFound
		}
		return nil, err
	}
	return wg, nil
}

func NewWordGroupRepoMongo(log logger.Logger, collection *mongo.Collection) *WordGroupRepoMongo {
	return &WordGroupRepoMongo{
		collection: collection,
		log:        log,
	}
}

func (wgr *WordGroupRepoMongo) FindByLangAndUser(ctx context.Context, langID string, userID string, def bool) (*domain.WordGroup, error) {
	userIDHEx, _ := primitive.ObjectIDFromHex(userID)
	langIDHEx, _ := primitive.ObjectIDFromHex(langID)
	wg := &domain.WordGroup{}
	result := wgr.collection.FindOne(ctx, bson.M{"lang_id": langIDHEx, "user_id": userIDHEx, "default": def})

	if err := result.Decode(wg); err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrNotFound
		}
		return nil, err
	}
	return wg, nil
}

func (wgr *WordGroupRepoMongo) Create(ctx context.Context, wordGroup *domain.WordGroup) (string, error) {
	wgm := &domain.WordGroupMongo{}
	err := copier.Copy(wgm, wordGroup)
	if err != nil {
		return "", err
	}
	userIDHEx, _ := primitive.ObjectIDFromHex(wordGroup.UserID)
	wgm.UserID = userIDHEx
	return create(ctx, wgr.collection, wgm)
}
