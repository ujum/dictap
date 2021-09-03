package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type WG struct {
	GroupID string    `bson:"group_id"`
	AddedAt time.Time `bson:"added_at"`
}

type Word struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
}

type WordGroupMongo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id"`
	Lang    primitive.ObjectID `bson:"lang_id"`
	Name    string             `bson:"name"`
	Default bool               `bson:"default"`
}

type WordGroup struct {
	ID      string `bson:"_id,omitempty"`
	UserID  string `bson:"user_id"`
	Name    string `bson:"name"`
	Lang    string `bson:"lang_id"`
	Default bool   `bson:"default"`
}
