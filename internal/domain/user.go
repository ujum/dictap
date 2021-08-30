package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Uid          string             `bson:"uid"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Phone        string             `bson:"phone"`
	Password     string             `bson:"password"`
	RegisteredAt time.Time          `bson:"registeredAt"`
}
