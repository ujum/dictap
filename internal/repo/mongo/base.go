package mongo

import (
	"context"
	derr "github.com/ujum/dictap/internal/error"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func create(ctx context.Context, collection *mongo.Collection, domainType interface{}) (string, error) {
	createdDomain, err := collection.InsertOne(ctx, domainType)
	if mongo.IsDuplicateKeyError(err) {
		return "", derr.ErrAlreadyExists
	}
	return createdDomain.InsertedID.(primitive.ObjectID).Hex(), nil
}
