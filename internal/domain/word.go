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
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
}

type WordGroupMongo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserUID string             `bson:"user_uid"`
	LangISO string             `bson:"lang_iso"`
	Name    string             `bson:"name"`
	Default bool               `bson:"default"`
}

type WordGroup struct {
	ID      string `bson:"_id,omitempty"`
	UserUID string `bson:"user_uid"`
	Name    string `bson:"name"`
	LangISO string `bson:"lang_iso"`
	Default bool   `bson:"default"`
}
