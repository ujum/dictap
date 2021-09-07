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
	UserUID string             `bson:"user_uid"`
	LangID  primitive.ObjectID `bson:"lang_id"`
	Name    string             `bson:"name"`
	Default bool               `bson:"default"`
}

type WordGroup struct {
	ID      string `bson:"_id,omitempty"`
	UserUID string `bson:"user_uid"`
	Name    string `bson:"name"`
	LangID  string `bson:"lang_id"`
	Default bool   `bson:"default"`
}
