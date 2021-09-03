package mongo

import (
	"context"
	"github.com/ujum/dictap/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func create(ctx context.Context, collection *mongo.Collection, domainType interface{}) (string, error) {
	createdDomain, err := collection.InsertOne(ctx, domainType)
	if mongo.IsDuplicateKeyError(err) {
		return "", domain.ErrAlreadyExists
	}
	return createdDomain.InsertedID.(primitive.ObjectID).Hex(), nil
}
