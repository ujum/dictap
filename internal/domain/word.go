package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type WG struct {
	GroupID primitive.ObjectID `bson:"group_id"`
	AddedAt time.Time          `bson:"added_at"`
}

type Word struct {
	ID      string    `bson:"_id,omitempty"`
	Name    string    `bson:"name"`
	AddedAt time.Time `bson:"added_at"`
}
