package mongo

import (
	"context"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type WordRepoMongo struct {
	log        logger.Logger
	collection *mongo.Collection
}

func NewWordRepoMongo(log logger.Logger, collection *mongo.Collection) *WordRepoMongo {
	return &WordRepoMongo{
		collection: collection,
		log:        log,
	}
}

func (wr *WordRepoMongo) AddToGroup(ctx context.Context, name string, groupID string) error {
	groupIDHEx, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = wr.collection.UpdateOne(ctx, bson.M{"name": name},
		bson.M{"$push": bson.M{"groups": bson.M{"group_id": groupIDHEx, "added_at": bsonx.Time(time.Now())}}})
	if err != nil {
		return err
	}
	return nil
}

func (wr *WordRepoMongo) FindByName(ctx context.Context, name string) (*domain.Word, error) {
	result := wr.collection.FindOne(ctx, bson.M{"name": name})
	word := &domain.Word{}

	if err := result.Decode(word); err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrNotFound
		}
		return nil, err
	}
	return word, nil
}

func (wr *WordRepoMongo) FindByGroup(ctx context.Context, groupID string) ([]*domain.Word, error) {
	groupIDHex, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}
	cursor, err := wr.collection.Find(ctx, bson.M{"groups.group_id": groupIDHex},
		&options.FindOptions{Sort: bson.D{{"groups.added_at", -1}}})

	var words []*domain.Word
	if err != nil {
		wr.log.Errorf("can't find words in group, reason: %v", err)
		return words, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &words)
	return words, err
}

func (wr *WordRepoMongo) FindByNameAndGroup(ctx context.Context, wordName string, groupID string) (*domain.Word, error) {
	groupIDHex, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}

	result := wr.collection.FindOne(ctx, bson.M{"name": wordName, "groups.group_id": groupIDHex})

	word := &domain.Word{}

	if err := result.Decode(word); err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrNotFound
		}
		return nil, err
	}
	return word, nil
}

func (wr *WordRepoMongo) Create(ctx context.Context, word *domain.Word) (string, error) {
	return create(ctx, wr.collection, word)
}

func (wr *WordRepoMongo) RemoveFromGroup(ctx context.Context, name string, groupID string) error {
	groupIDHEx, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = wr.collection.UpdateOne(ctx, bson.M{"name": name},
		bson.M{"$pull": bson.M{"groups": bson.M{"group_id": groupIDHEx}}})
	if err != nil {
		return err
	}
	return nil
}
