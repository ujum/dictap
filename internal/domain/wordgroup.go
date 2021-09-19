package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type WordGroupMongo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	UserUID     string             `bson:"user_uid"`
	LangBinding LangBinding        `bson:"lang_binding,omitempty"`
	Default     bool               `bson:"default"`
}

type WordGroup struct {
	ID          string      `bson:"_id,omitempty"`
	Name        string      `bson:"name"`
	UserUID     string      `bson:"user_uid"`
	LangBinding LangBinding `bson:"lang_binding,omitempty"`
	Default     bool        `bson:"default"`
}
