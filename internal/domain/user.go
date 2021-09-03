package domain

import (
	"time"
)

type User struct {
	ID           string    `bson:"_id,omitempty"`
	Uid          string    `bson:"uid"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Phone        string    `bson:"phone"`
	Password     string    `bson:"password"`
	RegisteredAt time.Time `bson:"registeredAt"`
}
