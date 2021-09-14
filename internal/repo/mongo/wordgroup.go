package mongo

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/ujum/dictap/internal/domain"
	derr "github.com/ujum/dictap/internal/error"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WordGroupRepoMongo struct {
	log        logger.Logger
	collection *mongo.Collection
}

func (wgr *WordGroupRepoMongo) FindAllByLangAndUser(ctx context.Context, langBinding *domain.LangBinding, userUID string) ([]*domain.WordGroup, error) {
	var wgs []*domain.WordGroup
	cursor, err := wgr.collection.Find(ctx, bson.M{"user_uid": userUID, "lang_binding": langBinding})
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

func (wgr *WordGroupRepoMongo) FindByIDAndUser(ctx context.Context, groupID string, userUID string) (*domain.WordGroup, error) {
	groupIDHEx, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}

	wg := &domain.WordGroup{}
	result := wgr.collection.FindOne(ctx, bson.M{"_id": groupIDHEx, "user_uid": userUID})
	if err := result.Decode(wg); err != nil {
		if err == mongo.ErrNoDocuments {
			err = derr.ErrNotFound
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

func (wgr *WordGroupRepoMongo) FindByLangAndUser(ctx context.Context, langBinding *domain.LangBinding, userUID string, def bool) (*domain.WordGroup, error) {
	wg := &domain.WordGroup{}
	result := wgr.collection.FindOne(ctx, bson.M{"lang_binding": langBinding, "user_uid": userUID, "default": def})

	if err := result.Decode(wg); err != nil {
		if err == mongo.ErrNoDocuments {
			err = derr.ErrNotFound
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
	if err != nil {
		return "", err
	}
	return create(ctx, wgr.collection, wgm)
}
